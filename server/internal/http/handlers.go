package http

import (
	"context"
	"encoding/json"
	"errors"
	"goQuiz/server/internal/cfg"
	"goQuiz/server/internal/clients"
	"goQuiz/server/internal/store"
	"goQuiz/server/internal/ws"
	wsHub "goQuiz/server/internal/ws"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

type Handler struct {
	Ref *store.Store
	Hub *wsHub.Hub
	Q   store.QuestionsRepo
	QM  *clients.Client
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

type CreateQuestionReq struct {
	Prompt       string   `json:"prompt"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correctIndex"`
	Category     *string  `json:"category,omitempty"`
	Difficulty   *string  `json:"difficulty,omitempty"`
}

func (h *Handler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req createRoomReq
	start := time.Now()

	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	var maxErr *http.MaxBytesError
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if errors.As(err, &maxErr) {
			http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
			if cfg.Debug {
				log.Printf("Handler -> CREATE | 413 payload too large {%dB}", maxErr.Limit)
			}
			return
		}
		http.Error(w, "invalid req body", http.StatusBadRequest)
		if cfg.Debug {
			log.Printf("Handler -> CREATE | err decode: %v", err)
		}
		return
	}

	if cfg.Debug {
		log.Printf("Handler -> CREATE | start name = {%s}", strings.TrimSpace(req.Name))
	}

	room := h.Ref.CreateRoom(req.Name)
	if cfg.Debug {
		log.Printf("Handler -> CREATE | ok code = {%s} , players = {%d} , dur = {%s}", room.Code, len(room.Players), time.Since(start))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)

}

func (h *Handler) GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	roomCode := chi.URLParam(r, "code")
	if cfg.Debug {
		log.Printf("Handler -> GET | extract code = {%s}", roomCode)
	}

	room, ok := h.Ref.GetRoom(roomCode)
	if !ok {
		http.Error(w, "invalid/room not found", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> GET | store get err")
		}
		return
	}

	if cfg.Debug {
		log.Printf("Handler -> GET | ok code = {%s} , players = {%d} , dur = {%s}", room.Code, len(room.Players), time.Since(start))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)

}

func (h *Handler) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	roomCode := chi.URLParam(r, "code")

	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	if cfg.Debug {
		log.Printf("Handler -> JOIN | extract code = {%s}", roomCode)
	}

	room, ok := h.Ref.GetRoom(roomCode)
	if !ok {
		http.Error(w, "invalid/room not found on join", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> JOIN | store get err")
		}
		return
	}

	var req joinReq
	var maxErr *http.MaxBytesError
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		if errors.As(err, &maxErr) {
			http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
			if cfg.Debug {
				log.Printf("Handler -> JOIN | 413 payload too large {%dB}", maxErr.Limit)
			}
			return
		}
		http.Error(w, "name required", http.StatusBadRequest)
		if cfg.Debug {
			log.Printf("Handler -> JOIN | err decode: %v", err)
		}
		return
	}

	if cfg.Debug {
		log.Printf("Handler -> JOIN | start code = {%s} , name = {%s}", roomCode, strings.TrimSpace(req.Name))
	}

	player, ok := h.Ref.JoinRoom(roomCode, req.Name)
	if !ok {
		http.Error(w, "cannot join room", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> JOIN | store join err")
		}
		return
	}
	room, _ = h.Ref.GetRoom(roomCode)

	if cfg.Debug {
		log.Printf("Handler -> JOIN | ok code = {%s} , id = {%s} , name = {%s}, dur = {%s}", roomCode, player.ID, player.Name, time.Since(start))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)

	lobby := lobby{Players: room.Players}
	bytes, err := json.Marshal(lobby)
	if err != nil {
		http.Error(w, "failed to marshal lobby data", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> JOIN | json byte marshal err")
		}
		return
	}

	h.Hub.Broadcast(roomCode, bytes)

}

func (h *Handler) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	roomCode := chi.URLParam(r, "code")
	if cfg.Debug {
		log.Printf("Handler -> LEAVE | extract code = {%s}", roomCode)
	}

	var req idReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cannot leave room", http.StatusBadRequest)
		if cfg.Debug {
			log.Printf("Handler -> LEAVE | err decode: %v", err)
		}
		return
	}

	if cfg.Debug {
		log.Printf("Handler -> LEAVE | start code = {%s} , id = {%s}", roomCode, req.ID)
	}

	room, ok := h.Ref.DropPlayer(roomCode, req.ID)
	if !ok {
		http.Error(w, "failed to drop player", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> LEAVE | store drop err")
		}
		return
	}

	if cfg.Debug {
		log.Printf("Handler -> LEAVE | ok code = {%s} , remPlayers = {%d} , dur = {%s}", roomCode, len(room.Players), time.Since(start))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)

	lobby := lobby{Players: room.Players}
	bytes, err := json.Marshal(lobby)
	if err != nil {
		http.Error(w, "failed to marshal lobby data/bum drop player", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> LEAVE | json byte marshal err")
		}
		return
	}

	h.Hub.Broadcast(roomCode, bytes)

}

func (h *Handler) WSHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	roomCode := chi.URLParam(r, "code")
	if cfg.Debug {
		log.Printf("Handler -> WS | extract code = {%s}", roomCode)
	}

	playerID := r.URL.Query().Get("playerId")
	if strings.TrimSpace(playerID) == "" {
		http.Error(w, "missing playerId", http.StatusBadRequest)
		if cfg.Debug {
			log.Printf("Handler -> WS | err bad playerID")
		}
		return
	}

	_, ok := h.Ref.GetRoom(roomCode)
	if !ok {
		http.Error(w, "room doesnt exist", http.StatusNotFound)
		if cfg.Debug {
			log.Printf("Handler -> WS | store get err")
		}
		return
	}

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil && cfg.Debug {
		log.Printf("Handler -> WS | err accept failed: %v", err)
		return
	}

	if cfg.Debug {
		log.Printf("Handler -> WS | ok code = {%s} , playerID = {%s} , dur = {%s}", roomCode, playerID, time.Since(start))
	}

	c := ws.NewConn(wsConn)
	h.Hub.Add(roomCode, playerID, c)

	ctx := r.Context()
	go c.WriteLoop(ctx)
	c.ReadLoop(ctx, h.Hub, roomCode, playerID)

	h.Hub.Remove(roomCode, playerID)

}

func (h *Handler) CreateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	if cfg.Debug {
		log.Printf("Handler -> CREATEQ | start")
	}

	var maxErr *http.MaxBytesError
	var req CreateQuestionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if errors.As(err, &maxErr) {
			http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
			if cfg.Debug {
				log.Printf("Handler -> CREATEQ | 413 payload too large {%dB}", maxErr.Limit)
			}
			return
		}
		http.Error(w, "invalid req body", http.StatusBadRequest)
		if cfg.Debug {
			log.Printf("Handler -> CREATEQ | err decode: %v", err)
		}
		return
	}

	p := strings.TrimSpace(req.Prompt)

	opts := make([]string, 0, len(req.Options))
	for _, o := range req.Options {
		o = strings.TrimSpace(o)
		if o != "" {
			opts = append(opts, o)
		}
	}

	var cat, diff *string
	if req.Category != nil && strings.TrimSpace(*req.Category) != "" {
		s := strings.TrimSpace(*req.Category)
		cat = &s
	}
	if req.Difficulty != nil && strings.TrimSpace(*req.Difficulty) != "" {
		s := strings.TrimSpace(*req.Difficulty)
		diff = &s
	}

	q := store.Question{
		Prompt:       p,
		Options:      opts,
		CorrectIndex: req.CorrectIndex,
		Category:     cat,
		Difficulty:   diff,
	}

	if err := q.Validate(); err != nil {
		http.Error(w, " invalid question: "+err.Error(), http.StatusBadRequest)
		if cfg.Debug {
			log.Printf("Handler -> CREATEQ | invalid question passed")
		}
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	id, err := h.Q.Insert(ctx, q)
	if err != nil {
		http.Error(w, "failed to insert question", http.StatusInternalServerError)
		if cfg.Debug {
			log.Printf("Handler -> CREATEQ | Insert reach, failed")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(map[string]any{"id": id})

	if cfg.Debug {
		log.Printf("Handler -> CREATEQ | ok dur = {%s}", time.Since(start))
	}

}
