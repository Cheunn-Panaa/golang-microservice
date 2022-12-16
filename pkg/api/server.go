package api

import (
	"cheunn-panaa/golang-microservice/config"
	"cheunn-panaa/golang-microservice/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

var (
	healthy int32
	ready   int32
)

type Server struct {
	router         *chi.Mux
	config         config.Config
	logger         *zap.Logger
	handler        http.Handler
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}

func NewServer(c config.Config, l *zap.Logger) (*Server, error) {

	var s = &Server{
		router: chi.NewRouter(),
		config: c,
		logger: l,
	}
	return s, nil
}

func (s *Server) ListenAndServe() (*http.Server, *int32, *int32) {
	ctx := context.Background()

	isReady := &atomic.Value{}
	isReady.Store(true)

	// Start the tracer
	s.initTracer(ctx)
	s.registerMiddlewares()

	s.initBaseRouter(isReady)

	s.handler = s.router
	srv := s.startServer()

	return srv, &healthy, &ready
}

func (s *Server) registerMiddlewares() {
	s.newOpenTelemetryMiddleware()
	s.router.Use(logger.RequestLogger(s.logger))
}

func (s *Server) startServer() *http.Server {

	// determine if the port is specified
	if s.config.Server.Port == 0 {
		// move on immediately
		return nil
	}
	address := fmt.Sprintf(":%d", s.config.Server.Port)

	srv := &http.Server{
		Addr:         address,
		Handler:      s.router,
		ReadTimeout:  s.config.Server.TimeoutRead,
		WriteTimeout: s.config.Server.TimeoutWrite,
		IdleTimeout:  s.config.Server.TimeoutIdle,
	}

	// start the server in the background
	go func() {
		s.logger.Sugar().Infof("Starting HTTP Server on port %d", s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	atomic.StoreInt32(&healthy, 1)
	atomic.StoreInt32(&ready, 1)

	return srv
}
