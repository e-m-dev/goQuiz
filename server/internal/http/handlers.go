package http

import (
	"encoding/json"
	"goQuiz/server/internal/store"
	"net/http"
)

type Handler struct {
	Ref *store.Store
}

type createRoomReq struct {
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
