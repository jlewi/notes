package main

import (
	"flag"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

const (
	loadAction = "/load"
)

type view string

const (
	errorView    view = "error"
	requestView  view = "request"
	responseView view = "response"
	rawView      view = "raw"

	loadErrorState = "/loadError"
	traceState     = "/trace"
)

var (
	port = flag.Int("port", 8080, "The port to listen on")
)

type mainLayout struct {
	app.Compo
	//main *markdownViewer
}

func (c *mainLayout) Render() app.UI {
	//log := zapr.NewLogger(zap.L())
	//if c.main == nil {
	//	//log.Info("Creating markdown viewer; I'm not sure why this is neccessary; shouldn't this be initialized in main")
	//	//c.main = &markdownViewer{}
	//}
	return app.Div().Class("main-layout").Body(
		app.Div().Class("header").Body(
			&tokenInput{},
		),
		//app.Div().Class("content").Body(
		//	app.Div().Class("sidebar").Body(
		//		&sideNavigator{},
		//	),
		//	app.Div().Class("main-window").Body(
		//		c.main,
		//	),
		//),
	)
}

// tokenInput is a component to enter an accessToken.
type tokenInput struct {
	app.Compo
	traceValue string
}

// The Render method is where the component appearance is defined.
func (h *tokenInput) Render() app.UI {
	return app.Div().Body(
		app.Input().
			Type("text").
			ID("inputBox").
			Value(""),
		app.Button().
			Text("Display").
			OnClick(func(ctx app.Context, e app.Event) {
				log := zapr.NewLogger(zap.L())
				accessToken := app.Window().GetElementByID("inputBox").Get("value").String()
				log.Info("Clicked", "accessToken", accessToken)
				if err := runGet(accessToken); err != nil {
					log.Error(err, "BigQuery request failed")
				}
			}),
	)
}

func runGet(accessToken string) error {
	log := zapr.NewLogger(zap.L())
	if !strings.HasPrefix(accessToken, "ya29") {
		log.Error(errors.New("Invalid access token"), "Access token doesn't start with ya29")
	}
	log.Info("Sending BigQuery get reuqest", "accessToken", accessToken)
	url := "https://bigquery.googleapis.com/bigquery/v2/projects/dev-sailplane/datasets/traces/tables/AgentTraces"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	//req.Header.Add("Origin", "http://localhost:8080")
	//req.Header.Add("Access-Control-Request-Method", "GET")
	//req.Header.Add("Access-Control-Request-Headers", "authorization")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to make request to %v", url)
	}

	fmt.Printf("Response:\n")
	fmt.Printf("StatusCode: %v\n", resp.StatusCode)
	fmt.Printf("Status: %v\n", resp.Status)
	//fmt.Printf("Headers:\n%+v\n", helpers.PrettyString(resp.Header))
	if resp.Body != nil {
		log.Info("Reading body")
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(err, "Failed to read response body")
		}

		log.Info("Read body", "body", string(b))
		fmt.Printf("\nBody: \n%v", string(b))
	} else {
		fmt.Print("No body")
	}

	// Is this the right way to verify CORS?
	if resp.StatusCode == http.StatusOK {
		fmt.Print("Request succeeded; CORS is probably supported")
	} else {
		fmt.Print("Request failed")
	}
	return nil
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

	// The first thing to do is to associate the tokenInput component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.Route("/", &mainLayout{})

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
			"/web/app.css", // Loads tokenInput.css file.
		},
	})

	log := zapr.NewLogger(zap.L())

	log.Info("Starting server", "port", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Error(err, "Failed to start server")
	}
}
