package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://192.168.0.100:5173"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/rooms", h.CreateRoomHandler)
	r.Get("/rooms/{code}", h.GetRoomHandler)
	r.Post("/rooms/{code}/join", h.JoinRoomHandler)
	r.Post("/rooms/{code}/leave", h.LeaveRoomHandler)
	r.Get("/ws/{code}", h.WSHandler)

	return r
}
