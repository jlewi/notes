package main

import (
	"flag"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"go.uber.org/zap"
	"net/http"
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

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
func (h *tokenInput) Render() app.UI {
	return app.Div().Body(
		app.Input().
			Type("text").
			ID("inputBox").
			Value("404d9f59c44936216791ed22a166a5e4"),
		app.Button().
			Text("Display").
			OnClick(func(ctx app.Context, e app.Event) {
				log := zapr.NewLogger(zap.L())
				accessToken := app.Window().GetElementByID("inputBox").Get("value").String()
				log.Info("Clicked", "accessToken", accessToken)
				//// Handle button click event here
				//client := pkg.AgentTracesClient{
				//	Endpoint: "http://localhost:8080",
				//}
				//traceID := app.Window().GetElementByID("inputBox").Get("value").String()
				//traceID = strings.TrimSpace(traceID)
				//if traceID == "" {
				//	h.traceValue = "No trace ID provided"
				//	h.Update()
				//	return
				//}
				//log := zapr.NewLogger(zap.L())
				//log.Info("Fetching trace", "traceID", traceID)
				//trace, err := client.GetTrace(ctx, traceID)
				//if err != nil {
				//	ctx.SetState(loadErrorState, err.Error())
				//	ctx.NewActionWithValue(loadAction, errorView)
				//} else {
				//	ctx.SetState(traceState, trace)
				//	ctx.NewActionWithValue(loadAction, rawView)
				//}
			}),
	)
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
		Description: "A viewer for Agent Traces",
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
