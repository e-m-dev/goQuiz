package http

import "github.com/go-chi/chi/v5"

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/rooms", h.CreateRoomHandler)
	r.Get("/rooms/{code}", h.GetRoomHandler)
	r.Post("/rooms/{code}/join", h.JoinRoomHandler)
	r.Post("/rooms/{code}/leave", h.LeaveRoomHandler)
	r.Get("/ws/{code}", h.WSHandler)
	return r
}
