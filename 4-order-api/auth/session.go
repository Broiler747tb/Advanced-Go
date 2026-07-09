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

// Create stores a new session for the given phone/code and returns its id.
func (s *SessionStore) Create(phone, code int) string {
	id := generateSessionId()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = Session{Phone: phone, Code: code, CreatedAt: time.Now()}
	return id
}

// Verify returns (phone, true) only if the session exists, the code matches,
// and it hasn't expired. The session is consumed (deleted) either way.
func (s *SessionStore) Verify(id string, code int) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[id]
	if !ok {
		return 0, false
	}
	delete(s.sessions, id) // one-time use
	if sess.Code != code || time.Since(sess.CreatedAt) > sessionTTL {
		return 0, false
	}
	return sess.Phone, true
}

func generateSessionId() string {
	b := make([]byte, 16)
	_, _ = crand.Read(b) // crypto/rand → unguessable session id
	return hex.EncodeToString(b)
}

// GenerateCode returns a random 4-digit confirmation code (1000-9999).
func GenerateCode() int {
	return rand.IntN(9000) + 1000
}
