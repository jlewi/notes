package main

import (
	"fmt"
	"os"

	"github.com/go-logr/zapr"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/jlewi/hydros/pkg/files"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/jlewi/tracesample/pkg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	util.SetupLogger("info", true)
	err := Run()
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}
}

func Run() error {
	log := zapr.NewLogger(zap.L())
	apiKeyUri := os.Getenv("HONEYCOMB_API_KEY_URI")
	if apiKeyUri == "" {
		return errors.New("HONEYCOMB_API_KEY_URI not set")
	}
	key, err := files.Read(apiKeyUri)
	if err != nil {
		return errors.Wrapf(err, "Could not read secret: %v", apiKeyUri)
	}
	log.Info("HONEYCOMB_API_KEY set enabling OTel", "uri", apiKeyUri)
	// Enable multi-span attributes
	bsp := honeycomb.NewBaggageSpanProcessor()

	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithSpanProcessor(bsp),
		honeycomb.WithApiKey(string(key)),
	)
	if err != nil {
		return errors.Wrapf(err, "error setting up OTel SDK")
	}
	defer otelShutdown()

	httpServer := pkg.HttpServer{
		Port: 8080,
	}

	if err := httpServer.Setup(); err != nil {
		return errors.Wrapf(err, "error setting up http server")
	}

	go func() {
		httpServer.StartAndBlock()
	}()

	pkg.RunServer(50051)

	// Loop forever
	select {} // block forever
}
