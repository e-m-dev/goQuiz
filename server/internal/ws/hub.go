package ws

import (
	"context"
	"goQuiz/server/internal/cfg"
	"log"
	"sync"

	"nhooyr.io/websocket"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[string]*Conn
}

type Conn struct {
	ws   *websocket.Conn
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[string]*Conn),
	}
}

func (h *Hub) Add(roomCode string, playerID string, c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.rooms[roomCode]; !exists {
		h.rooms[roomCode] = make(map[string]*Conn)
	}

	h.rooms[roomCode][playerID] = c

	if cfg.Debug {
		log.Printf("Hub -> ADD | room = {%s} , id = {%s} , total = {%d}", roomCode, playerID, len(h.rooms[roomCode]))
	}
}

func (h *Hub) Remove(roomCode string, playerID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if players, exists := h.rooms[roomCode]; exists {
		delete(players, playerID)
		if len(players) == 0 {
			delete(h.rooms, roomCode)
		}
	}

	if cfg.Debug {
		log.Printf("Hub -> REMOVE | room = {%s} , id = {%s} , total = {%d}", roomCode, playerID, len(h.rooms[roomCode]))
	}
}

func (h *Hub) Broadcast(roomCode string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if players, exists := h.rooms[roomCode]; exists && cfg.Debug {
		log.Printf("Hub -> BCAST | room = {%s} , receipients = {%d} , bytes = {%d}", roomCode, len(players), len(payload))
		for _, conn := range players {
			select {
			case conn.send <- payload:
			default:
			}
		}
	}
}

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{
		ws:   ws,
		send: make(chan []byte, 256),
	}
}

func (c *Conn) WriteLoop(ctx context.Context) {
	defer func() {
		_ = c.ws.Close(websocket.StatusNormalClosure, "writer closed")
		if cfg.Debug {
			log.Printf("Hub -> WRITE | close err = {%s}", ctx.Err())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			if err := c.ws.Write(ctx, websocket.MessageText, msg); err != nil {
				return
			}
		}
	}
}

func (c *Conn) ReadLoop(ctx context.Context, hub *Hub, roomCode string, playerID string) {
	defer func() {
		close(c.send)
		if cfg.Debug {
			log.Printf("Hub -> READ | read close room = {%s} , id = {%s} , err = {%v}", roomCode, playerID, ctx.Err())
		}
	}()
	defer hub.Remove(roomCode, playerID)
	for {
		typ, _, err := c.ws.Read(ctx)
		if err != nil {
			return
		}

		if typ == websocket.MessageText {
			//hub.Broadcast(roomCode, msg)
		}
	}
}
