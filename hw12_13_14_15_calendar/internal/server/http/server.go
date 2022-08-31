package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	config *config.HTTPConf
	logger app.Logger
}

type Application interface { // TODO
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}

func NewServer(logger app.Logger, app Application, config *config.HTTPConf) *Server {
	return &Server{
		logger: logger,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.logger.Info("Start listen at " + addr)
	router := http.NewServeMux()
	router.HandleFunc("/", hello)
	configuredRouter := loggingMiddleware(s.logger)(router)
	err := http.ListenAndServe(addr, configuredRouter)
	if err != nil {
		s.logger.Error("Failed listen at " + addr)
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
