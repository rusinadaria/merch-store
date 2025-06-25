package handlers

import (
	"merch-store/internal/services"
	"net/http"
	"github.com/go-chi/chi"
	"log/slog"
	"merch-store/internal/handlers/middleware"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.LoggerMiddlewareWrapper(logger))

	mw := middleware.Handler{}

	r.Group(func(r chi.Router) {
		r.Post("/api/auth", h.AddUserHandler)

		r.With(mw.AuthMiddleware).Group(func(r chi.Router) {
			r.Get("/api/info", h.InfoHandler)
			r.Get("/api/buy/{item}", h.BuyItemHandler)
			r.Post("/api/sendCoin", h.SendHandler)
		})
	})
	return r
}