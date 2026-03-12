package session

import (
	"sync"
	"time"

	dbsqlc "inventari/api/internal/db/sqlc"
)

const ttl = 24 * time.Hour

type entry struct {
	user      dbsqlc.AppUser
	expiresAt time.Time
}

// Store is a thread-safe in-memory session store.
// Sessions are volatile — they are lost on server restart.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]entry
}

func NewStore() *Store {
	s := &Store{sessions: make(map[string]entry)}
	go s.cleanupLoop()
	return s
}

func (s *Store) Set(token string, user dbsqlc.AppUser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = entry{user: user, expiresAt: time.Now().Add(ttl)}
}

func (s *Store) Get(token string) (dbsqlc.AppUser, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.sessions[token]
	if !ok || time.Now().After(e.expiresAt) {
		return dbsqlc.AppUser{}, false
	}
	return e.user, true
}

func (s *Store) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

// cleanupLoop removes expired sessions every 15 minutes.
func (s *Store) cleanupLoop() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.evict()
	}
}

func (s *Store) evict() {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for token, e := range s.sessions {
		if now.After(e.expiresAt) {
			delete(s.sessions, token)
		}
	}
}
