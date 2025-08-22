package store

import (
	"math/rand/v2"
	"strings"
	"sync"
)

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host bool   `json:"host"`
}

type Room struct {
	Code    string `json:"code"`
	Name    string
	Players []Player `json:"players"`
}

type Store struct {
	rooms map[string]Room
	mu    sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		rooms: make(map[string]Room),
	}
}

func (s *Store) CreateRoom(name string) Room {
	roomCode := randomCode(6)

	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		if _, exists := s.rooms[roomCode]; !exists {
			break
		}
		roomCode = randomCode(6)
	}

	room := Room{Code: roomCode, Name: name, Players: make([]Player, 0)}
	s.rooms[roomCode] = room
	return room
}

func (s *Store) GetRoom(code string) (Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[code]
	return room, exists
}

func (s *Store) JoinRoom(code string, playerName string) (Player, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms[code]
	if !ok {
		return Player{}, false
	}

	name := strings.TrimSpace(playerName)
	if name == "" {
		return Player{}, false
	}

	for _, p := range room.Players {
		if strings.EqualFold(p.Name, name) {
			return Player{}, false
		}
	}

	isHost := false
	if len(room.Players) <= 0 {
		isHost = true
	}

	player := Player{ID: randomCode(16), Name: name, Host: isHost}
	room.Players = append(room.Players, player)

	s.rooms[code] = room

	return player, true

}

func (s *Store) DropPlayer(code string, id string) (Room, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms[code]
	if !ok {
		return Room{}, false
	}

	for i := range room.Players {

		if room.Players[i].ID == id {
			wasHost := room.Players[i].Host

			room.Players = append(room.Players[:i], room.Players[i+1:]...)

			if len(room.Players) <= 0 {
				delete(s.rooms, code)
				return Room{}, true
			}

			if wasHost {
				room.Players[0].Host = true
			}

			s.rooms[code] = room
			return room, true
		}

	}

	s.rooms[code] = room
	return room, false
}

func randomCode(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}
