package api

import (
	"net/http"
	"sync/atomic"
)

// Healthz godoc
// @Summary Liveness check
// @Description used by Kubernetes liveness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /healthz [get]
// @Success 200 {string} string "OK"
func (s *Server) Healthz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Readyz godoc
// @Summary Readiness check
// @Description used by Kubernetes readiness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz [get]
// @Success 200 {string} string "OK"
func (s *Server) Readyz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&ready) == 1 {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
