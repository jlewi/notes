package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"strings"
)

const (
	updateAction          = "/updateOutput"
	tokenState            = "accessToken"
	useClientLibState     = "useClientLib"
	addCustomHeadersState = "useAddCustomHeaders"
	cloudPlatformScope    = "https://www.googleapis.com/auth/cloud-platform"
)

var (
	port = flag.Int("port", 8080, "The port to listen on")
)

// mainApp is the app to test CORS.
type mainApp struct {
	app.Compo
	traceValue string
}

// The Render method is where the component appearance is defined.
func (h *mainApp) Render() app.UI {
	return app.Div().Class("main-layout").Body(
		&InputBoxes{
			projectValue: "dev-sailplane",
			datasetValue: "traces",
			tableValue:   "AgentTraces",
		},
		&OptionsComponent{
			useClientLib:     false,
			addCustomHeaders: false,
		},
		app.Button().
			Text("Run").
			OnClick(func(ctx app.Context, e app.Event) {
				log := zapr.NewLogger(zap.L())
				if err := h.runGet(ctx); err != nil {
					log.Error(err, "BigQuery request failed")
				}
			}),
		&OutputBox{},
	)
}

// OnTokenChange handles the change event
func (h *mainApp) OnTokenChange(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("value").String()
	log := zapr.NewLogger(zap.L())
	log.Info("Setting accessToken")
	ctx.SetState(tokenState, value)
}

// InputBoxes is a component that includes three input boxes.
type InputBoxes struct {
	app.Compo
	projectValue string
	datasetValue string
	tableValue   string
}

func (h *InputBoxes) OnMount(ctx app.Context) {
	// Initialize the context with whatever the current values are
	log := zapr.NewLogger(zap.L())
	log.Info("initializing state variables", "table", h.tableValue, "project", h.projectValue, "dataset", h.datasetValue)
	ctx.SetState("table", h.tableValue)
	ctx.SetState("project", h.projectValue)
	ctx.SetState("dataset", h.datasetValue)
}

// The Render method is where the component appearance is defined.
func (h *InputBoxes) Render() app.UI {
	return app.Div().Body(
		app.P().Text("Project:"),
		app.Input().
			Type("text").
			Value(h.projectValue).
			OnChange(h.OnProjectChange),
		app.P().Text("Dataset:"),
		app.Input().
			Type("text").
			Value(h.datasetValue).OnChange(h.OnDatasetChange),
		app.P().Text("Table:"),
		app.Input().
			Type("text").
			Value(h.tableValue).OnChange(h.OnTableChange),
	)
}

// OnProjectChange handles the change event of the project input box.
func (h *InputBoxes) OnProjectChange(ctx app.Context, e app.Event) {
	h.projectValue = ctx.JSSrc().Get("value").String()
	log := zapr.NewLogger(zap.L())
	log.Info("Updating project", "newValue", h.projectValue)
	ctx.SetState("project", h.projectValue)
}

// OnDatasetChange handles the change event of the dataset input box.
func (h *InputBoxes) OnDatasetChange(ctx app.Context, e app.Event) {
	h.datasetValue = ctx.JSSrc().Get("value").String()
	ctx.SetState("dataset", h.datasetValue)
}

// OnTableChange handles the change event of the table input box.
func (h *InputBoxes) OnTableChange(ctx app.Context, e app.Event) {
	h.tableValue = ctx.JSSrc().Get("value").String()
	ctx.SetState("table", h.tableValue)
}

func (h *mainApp) fetchAccessToken(ctx app.Context) error {
	endpoint := app.Window().URL().String() + "/getaccesstoken"

	resp, err := http.Get(endpoint)
	log := zapr.NewLogger(zap.L())
	if err != nil {
		log.Error(err, "Failed to get access token")
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err := errors.Errorf("Failed to obtain access token; StatusCode %v", resp.StatusCode)
		log.Error(err, "Failed to obtain access token", "statusCode", resp.StatusCode)
		return err
	}
	defer resp.Body.Close()

	// Read the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err, "Failed to read access token")
		return err
	}

	ctx.SetState(tokenState, string(data))
	return nil
}

func (h *mainApp) runGet(ctx app.Context) error {
	var project string
	var dataset string
	var table string
	var accessToken string
	ctx.GetState(tokenState, &accessToken)
	ctx.GetState("project", &project)
	ctx.GetState("dataset", &dataset)
	ctx.GetState("table", &table)
	if accessToken == "" {
		if err := h.fetchAccessToken(ctx); err != nil {
			return err
		}

		ctx.GetState(tokenState, &accessToken)

		if accessToken == "" {
			return errors.New("AccessToken not set in context after fetch")
		}
	}

	var useClientLib bool
	var addCustomHeaders bool
	ctx.GetState(useClientLibState, &useClientLib)
	ctx.GetState(addCustomHeadersState, &addCustomHeaders)

	sb := &strings.Builder{}

	if useClientLib {
		result, err := h.runGetWithClient(project, dataset, table, accessToken)

		if err != nil {
			sb.WriteString(err.Error())
		} else {
			sb.WriteString(result)
		}
	} else {
		result, err := h.runGetWithHttp(project, dataset, table, accessToken, addCustomHeaders)

		if err != nil {
			sb.WriteString(err.Error())
		} else {
			sb.WriteString(result)
		}
	}

	ctx.NewActionWithValue(updateAction, sb.String())
	return nil
}

// perform the get using the bigquery client library
func (h *mainApp) runGetWithClient(project string, dataset string, table string, accessToken string) (string, error) {
	log := zapr.NewLogger(zap.L())
	if accessToken == "" {
		return "", errors.New("accessToken is empty")
	}
	log.Info("Using runGetWithClient")
	// Create an OAuth2 token source using the access token
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	client, err := bigquery.NewClient(context.Background(), project, option.WithTokenSource(tokenSource))
	if err != nil {
		return "", errors.Wrapf(err, "Failed to create bigquery client")
	}
	defer client.Close()

	md, err := client.DatasetInProject(project, dataset).Table(table).Metadata(context.Background())
	if err != nil {
		return "", errors.Wrapf(err, "Failed to fetch table")
	}
	jsonData, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return "", errors.Wrapf(err, "Failed to marshal BigQuery output to JSON")
	}

	return string(jsonData), nil
}

func (h *mainApp) runGetWithHttp(project string, dataset string, table string, accessToken string, addCustomHeaders bool) (string, error) {
	log := zapr.NewLogger(zap.L())
	sb := &strings.Builder{}
	if !strings.HasPrefix(accessToken, "ya29") {
		log.Error(errors.New("Invalid access token"), "Access token doesn't start with ya29")
	}
	log.Info("Sending BigQuery get request using http")
	url := fmt.Sprintf("https://bigquery.googleapis.com/bigquery/v2/projects/%s/datasets/%s/tables/%s", project, dataset, table)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	if addCustomHeaders {
		log.Info("Adding custom headers")
		req.Header.Add("X-Goog-Api-Client", "1234")
		req.Header.Add("X-Cloud-Trace-Context", "1234")
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to make request to %v", url)
	}

	fmt.Fprintf(sb, "Response:\n")
	fmt.Fprintf(sb, "StatusCode: %v\n", resp.StatusCode)
	fmt.Fprintf(sb, "Status: %v\n", resp.Status)

	if resp.Body != nil {
		log.Info("Reading body")
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrapf(err, "Failed to read response body")
		}

		log.Info("Read body", "body", string(b))
		fmt.Fprintf(sb, "\nBody: \n%v", string(b))
	} else {
		fmt.Fprint(sb, "No body")
	}

	// Is this the right way to verify CORS?
	if resp.StatusCode == http.StatusOK {
		fmt.Fprint(sb, "Request succeeded; CORS is probably supported")
	} else {
		fmt.Fprint(sb, "Request failed")
	}
	return sb.String(), nil
}

// OutputBox is where the output is displayed.
type OutputBox struct {
	app.Compo
	outputValue string
}

func (h *OutputBox) Render() app.UI {
	return app.Div().Class("main-window").Body(
		app.Textarea().Class("textarea").
			Text(h.outputValue).
			ReadOnly(true))
}

func (m *OutputBox) OnMount(ctx app.Context) {
	// Registering action handler for the updateAction.
	ctx.Handle(updateAction, m.handleUpdateAction)
}

func (m *OutputBox) handleUpdateAction(ctx app.Context, action app.Action) {
	log := zapr.NewLogger(zap.L())
	output, ok := action.Value.(string) // Checks if a name was given.
	if !ok {
		log.Error(errors.New("No output provided"), "Invalid action")
		return
	}
	m.outputValue = output
	// Trigger the render method to update the view
	m.Update()
}

func getAccessToken() (string, error) {
	creds, err := google.DefaultTokenSource(context.Background(), cloudPlatformScope)
	if err != nil {
		return "", err
	}

	// Get an OAuth token
	token, err := creds.Token()
	if err != nil {
		fmt.Println("Failed to get an OAuth token:", err)
		return "", err
	}

	// Print the OAuth token
	fmt.Println("OAuth token:", token.AccessToken)
	return token.AccessToken, nil
}

// TODO(jeremy): We should really return the token as a Token object (serialized as JSON)
// so we transmit the expiration time to the client.

func getAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getAccessToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(token))
}

// OptionsComponent collects input about how to make the BigQuery request.
type OptionsComponent struct {
	app.Compo
	useClientLib     bool
	addCustomHeaders bool
}

func (c *OptionsComponent) Render() app.UI {
	return app.Div().Body(
		app.Input().
			Type("checkbox").
			ID("checkbox1").
			Checked(c.useClientLib).
			OnChange(c.OnuseClientLibChange),
		app.Label().For("checkbox1").Text("Use GCP Client Library"),
		app.Br(),
		app.Input().
			Type("checkbox").
			ID("checkbox2").
			Checked(c.addCustomHeaders).
			OnChange(c.OnAddCustomHeadersChange),
		app.Label().For("checkbox2").Text("Add Custom Headers"),
	)
}

func (c *OptionsComponent) OnuseClientLibChange(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("checked").Bool()
	ctx.SetState(useClientLibState, value)
}

func (c *OptionsComponent) OnAddCustomHeadersChange(ctx app.Context, e app.Event) {
	value := ctx.JSSrc().Get("checked").Bool()
	ctx.SetState(addCustomHeadersState, value)
}

func (c *OptionsComponent) OnMount(ctx app.Context) {
	// Initialize the context with whatever the current values are
	ctx.SetState(useClientLibState, c.useClientLib)
	ctx.SetState(addCustomHeadersState, c.addCustomHeaders)
}

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {
	// How does logging work in go-app?
	c := zap.NewProductionConfig()

	// Use the keys used by cloud logging
	// https://cloud.google.com/logging/docs/structured-logging
	c.EncoderConfig.LevelKey = "severity"
	c.EncoderConfig.TimeKey = "time"
	c.EncoderConfig.MessageKey = "message"

	c.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	newLogger, err := c.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize zap logger; error %v", err))
	}

	zap.ReplaceGlobals(newLogger)

	// The first thing to do is to associate the mainApp component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &mainApp{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// When executed on the server-side, RunWhenOnBrowser() does nothing, which
	// lets room for server implementation without the need for precompiling
	// instructions.
	app.RunWhenOnBrowser()

	// Put instructions that should run on the server here
	flag.Parse()

	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "AgentTraceView",
		Description: "BigQueryViewer",
		Styles: []string{
			"/web/app.css", // Loads mainApp.css file.
		},
	})

	// Add a server side method to get credentials
	http.HandleFunc("/getaccesstoken", getAccessTokenHandler)

	log := zapr.NewLogger(zap.L())

	log.Info("Starting server", "port", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Error(err, "Failed to start server")
	}
}
