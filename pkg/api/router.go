package api

import (
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) initBaseRouter(ready *atomic.Value) *chi.Mux {
	s.logger.Debug("Initializing basic handler/router")

	r := s.router
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Get("/healthz", s.Healthz)
	r.Get("/readyz", s.Readyz)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"), //The url pointing to API definition
	))

	return r
}
