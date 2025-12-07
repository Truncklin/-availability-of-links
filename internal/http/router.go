package http

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"AvailabilityLinks/internal/handlers"
)

func InitRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handlers.Health)
		r.Post("/links", h.SubmitLinks)
		r.Get("/report", h.GetReport)
	})
	return r
}
