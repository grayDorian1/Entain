package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/grayDorian1/Entain/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handler) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/health", h.HealthCheck)
	r.Post("/user/{userId}/transaction", h.HandleTransaction)
	r.Get("/user/{userId}/balance", h.HandleBalance)

	return r
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}