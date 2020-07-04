package anansi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// DefaultMiddleware sets up some middlware for a router as well as
// liveliness(/) and not found handlers:
// - RequestID
// - RealIP
// - Recoverer
// - RedirectSlashes
// - Compress(with compression level of 5)
// - Timeout(with 1 minute)
func DefaultMiddleware(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Timeout(time.Minute))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Up and Running!"))
		if err != nil {
			panic(err)
		}
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Whoops!! This route doesn't exist", 404)
	})
}

// CORS sets CORS for the handler. It enables localhost by default when appEnv
// is not "dev".
func CORS(router *chi.Mux, appEnv string, origins ...string) {
	if appEnv == "dev" {
		router.Use(DevCORS().Handler)
	} else {
		origins = append(origins, "http://localhost")
		router.Use(SecureCORS(origins...).Handler)
	}
}
