package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type Server struct {
	httpServer *http.Server
	logger     Logger
}

func NewServer(logger Logger, host, port string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    host + ":" + port,
			Handler: handler,

			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting http server...", map[string]interface{}{"address": s.httpServer.Addr})
	err := s.httpServer.ListenAndServe()
	if err != nil {
		s.logger.Error("error while starting http server", map[string]interface{}{"error": err})
	}

	return err
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
