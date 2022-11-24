package api

import (
	"net/http"
)

// Healthz godoc
// @Summary Liveness check
// @Description used by Kubernetes liveness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /healthz [get]
// @Success 200 {string} string "OK"
func Healthz(w http.ResponseWriter, _ *http.Request) {
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
func Readyz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	//if atomic.LoadInt32(&ready) == 1 {
	//	w.WriteHeader(http.StatusOK)
	//	return
	//}
	//w.WriteHeader(http.StatusServiceUnavailable)
}
