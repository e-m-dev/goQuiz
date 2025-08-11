package http

import (
	"encoding/json"
	"goQuiz/server/internal/store"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Ref *store.Store
}

type createRoomReq struct {
	Name string `json:"name"`
}

type joinReq struct {
	Name string `json:"name"`
}

func (h *Handler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req createRoomReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid req body", http.StatusBadRequest)
		return
	}

	room := h.Ref.CreateRoom(req.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)

}

func (h *Handler) GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := chi.URLParam(r, "code")

	room, ok := h.Ref.GetRoom(roomCode)
	if !ok {
		http.Error(w, "invalid/room not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)

}

func (h *Handler) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := chi.URLParam(r, "code")

	_, ok := h.Ref.GetRoom(roomCode)
	if !ok {
		http.Error(w, "invalid/room not found on join", http.StatusNotFound)
		return
	}

	var req joinReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	player, ok := h.Ref.JoinRoom(roomCode, req.Name)
	if !ok {
		http.Error(w, "cannot join room", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)

}
