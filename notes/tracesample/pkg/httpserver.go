package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"
	"github.com/jlewi/p22h/backend/api"
	"github.com/jlewi/p22h/backend/pkg/debug"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"
)

type HttpServer struct {
	log    logr.Logger
	Port   int
	router *mux.Router
	srv    *http.Server
}

func (s *HttpServer) Setup() error {
	s.log = zapr.NewLogger(zap.L())
	return s.addRoutes()
}

// StartAndBlock starts the server and blocks.
func (s *HttpServer) StartAndBlock() {
	log := s.log
	log.Info("Binding all network interfaces", "port", s.Port)
	s.srv = &http.Server{Addr: fmt.Sprintf(":%v", s.Port), Handler: s.router}

	s.trapInterrupt()
	err := s.srv.ListenAndServe()

	if err != nil {
		if err == http.ErrServerClosed {
			log.Info("HttpServer has been shutdown")
		} else {
			log.Error(err, "Server aborted with error")
		}
	}
}

func (s *HttpServer) addRoutes() error {
	log := zapr.NewLogger(zap.L())
	router := mux.NewRouter().StrictSlash(true)
	// TODO(jeremy): Should we get the service name from an environemnt variable or something?

	if true {
		router.Use(otelmux.Middleware("foyle"))
	} else {
		log.Info("Warning; not adding opentelemetry middleware to http server")
	}
	s.router = router

	hPath := "/healthz"
	log.Info("Registering health check", "path", hPath)
	router.HandleFunc(hPath, s.healthCheck)

	router.NotFoundHandler = http.HandlerFunc(s.notFoundHandler)

	return nil
}

func (s *HttpServer) writeStatus(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := api.RequestStatus{
		Kind:    "RequestStatus",
		Message: message,
		Code:    code,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		s.log.Error(err, "Failed to marshal RequestStatus", "RequestStatus", resp, "code", code)
	}

	if code != http.StatusOK {
		caller := debug.ThisCaller()
		s.log.Info("HTTP error", "RequestStatus", resp, "code", code, "caller", caller)
	}
}
func (s *HttpServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	s.writeStatus(w, "HttpServer server is running", http.StatusOK)
}

func (s *HttpServer) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.writeStatus(w, fmt.Sprintf("HttpServer server doesn't handle the path; url: %v", r.URL), http.StatusNotFound)
}

// trapInterrupt waits for a shutdown signal and shutsdown the server
func (s *HttpServer) trapInterrupt() {
	sigs := make(chan os.Signal, 10)
	// SIGSTOP and SIGTERM can't be caught; however SIGINT works as expected when using ctl-z
	// to interrupt the process
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		msg := <-sigs
		s.log.Info("Received shutdown signal", "sig", msg)
		if err := s.srv.Shutdown(context.Background()); err != nil {
			s.log.Error(err, "Error shutting down server.")
		}
	}()
}
