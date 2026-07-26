package auth

import (
	"testing"
	"time"
)

func TestSessionStore_OneTimeUse(t *testing.T) {
	store := NewSessionStore()
	id := store.Create(89990009900, 1234)

	if _, ok := store.Verify(id, 1234); !ok {
		t.Fatal("first verify with correct code should succeed")
	}

	if _, ok := store.Verify(id, 1234); ok {
		t.Fatal("second verify of the same session should fail (one-time use)")
	}
}

func TestSessionStore_WrongCodeConsumesSession(t *testing.T) {
	store := NewSessionStore()
	id := store.Create(89990009900, 1234)

	if _, ok := store.Verify(id, 9999); ok {
		t.Fatal("verify with wrong code should fail")
	}

	if _, ok := store.Verify(id, 1234); ok {
		t.Fatal("session should be consumed even after a wrong-code attempt")
	}
}

func TestSessionStore_Expired(t *testing.T) {
	store := NewSessionStore()
	id := store.Create(89990009900, 1234)

	store.mu.Lock()
	sess := store.sessions[id]
	sess.CreatedAt = time.Now().Add(-sessionTTL - time.Minute)
	store.sessions[id] = sess
	store.mu.Unlock()

	if _, ok := store.Verify(id, 1234); ok {
		t.Fatal("expired session should not verify even with the correct code")
	}
}

func TestGenerateCode_InRange(t *testing.T) {
	for i := 0; i < 1000; i++ {
		code := GenerateCode()
		if code < 1000 || code > 9999 {
			t.Fatalf("GenerateCode() = %d, want a 4-digit number in [1000,9999]", code)
		}
	}
}
