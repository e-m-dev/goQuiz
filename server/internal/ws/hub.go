package ws

import (
	"context"
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
}

func (h *Hub) Broadcast(roomCode string, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if players, exists := h.rooms[roomCode]; exists {
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
	defer func() { _ = c.ws.Close(websocket.StatusNormalClosure, "writer closed") }()
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
	defer close(c.send)
	for {
		typ, msg, err := c.ws.Read(ctx)
		if err != nil {
			return
		}

		if typ == websocket.MessageText {
			hub.Broadcast(roomCode, msg)
		}
	}
}
