package http

import (
	"encoding/json"
	"goQuiz/server/internal/store"
	"goQuiz/server/internal/ws"
	wsHub "goQuiz/server/internal/ws"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

type Handler struct {
	Ref *store.Store
	Hub *wsHub.Hub
}

type createRoomReq struct {
	Name string `json:"name"`
}

type joinReq struct {
	Name string `json:"name"`
}

type idReq struct {
	ID string `json:"id"`
}

type lobby struct {
	Players []store.Player `json:"players"`
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

	room, ok := h.Ref.GetRoom(roomCode)
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

	lobby := lobby{Players: room.Players}
	bytes, err := json.Marshal(lobby)
	if err != nil {
		http.Error(w, "failed to marshal lobby data", http.StatusNotFound)
		return
	}

	h.Hub.Broadcast(roomCode, bytes)

}

func (h *Handler) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := chi.URLParam(r, "code")

	var req idReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot leave room", http.StatusBadRequest)
		return
	}

	room, ok := h.Ref.DropPlayer(roomCode, req.ID)
	if !ok {
		http.Error(w, "failed to drop player", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)

	lobby := lobby{Players: room.Players}
	bytes, err := json.Marshal(lobby)
	if err != nil {
		http.Error(w, "failed to marshal lobby data/bum drop player", http.StatusNotFound)
		return
	}

	h.Hub.Broadcast(roomCode, bytes)

}

func (h *Handler) WSHandler(w http.ResponseWriter, r *http.Request) {
	roomCode := chi.URLParam(r, "code")

	playerID := r.URL.Query().Get("playerId")
	if strings.TrimSpace(playerID) == "" {
		http.Error(w, "missing playerId", http.StatusBadRequest)
		return
	}

	_, ok := h.Ref.GetRoom(roomCode)
	if !ok {
		http.Error(w, "room doesnt exist", http.StatusNotFound)
		return
	}

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Printf("ws accept failed: &v", err)
		return
	}

	c := ws.NewConn(wsConn)
	h.Hub.Add(roomCode, playerID, c)

	ctx := r.Context()
	go c.WriteLoop(ctx)
	c.ReadLoop(ctx, h.Hub, roomCode, playerID)

	h.Hub.Remove(roomCode, playerID)

}
