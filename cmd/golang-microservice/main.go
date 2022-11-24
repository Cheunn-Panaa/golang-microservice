package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/cheunn-panaa/golang-microservice/pkg/api"
)

func main() {

	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Printf("Readyz probe is negative by default...")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Printf("Readyz probe is positive.")
	}()

	router := initRouter(isReady)
	http.ListenAndServe(":3000", router)
}

func initRouter(ready *atomic.Value) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/healthz", api.Healthz)

	r.Get("/readyz", api.Readyz(ready))

	return r
}
