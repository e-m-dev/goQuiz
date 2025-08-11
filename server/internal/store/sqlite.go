package store

import (
	"math/rand/v2"
	"sync"
)

type Room struct {
	Code string `json:"code"`
	Name string
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

	room := Room{Code: roomCode, Name: name}
	s.rooms[roomCode] = room
	return room
}

func (s *Store) GetRoom(code string) (Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[code]
	return room, exists
}

func randomCode(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}
