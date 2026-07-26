package auth

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand/v2"
	"sync"
	"time"
)

const sessionTTL = 5 * time.Minute

type Session struct {
	Phone     int
	Code      int
	CreatedAt time.Time
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: map[string]Session{}}
}

func (s *SessionStore) Create(phone, code int) string {
	id := generateSessionId()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = Session{Phone: phone, Code: code, CreatedAt: time.Now()}
	return id
}

func (s *SessionStore) Verify(id string, code int) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[id]
	if !ok {
		return 0, false
	}
	delete(s.sessions, id)
	if sess.Code != code || time.Since(sess.CreatedAt) > sessionTTL {
		return 0, false
	}
	return sess.Phone, true
}

func generateSessionId() string {
	b := make([]byte, 16)
	_, _ = crand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateCode() int {
	return rand.IntN(9000) + 1000
}
