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

func (c *Conn) WriteLoop() {
	for msg := range c.send {
		if err := c.ws.Write(context.Background(), websocket.MessageText, msg); err != nil {
			break
		}
	}
	_ = c.ws.Close(websocket.StatusNormalClosure, "")
}

func (c *Conn) ReadLoop(hub *Hub, roomCode string, playerID string) {
	for {
		typ, msg, err := c.ws.Read(context.Background())
		if err != nil {
			break
		}

		if typ != websocket.MessageText {
			continue
		}
		hub.Broadcast(roomCode, msg)
	}
	_ = c.ws.Close(websocket.StatusNormalClosure, "")
	close(c.send)
}
