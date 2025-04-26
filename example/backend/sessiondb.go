package main

import (
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

// SessionDB provides thread-safe in-memory storage for WebAuthn session data.
type SessionDB struct {
	sessions map[string]*webauthn.SessionData
	mu       sync.RWMutex
}

// NewSessionDB creates a new instance of SessionDB with an empty sessions map.
func NewSessionDB() *SessionDB {
	return &SessionDB{
		sessions: make(map[string]*webauthn.SessionData),
	}
}

// SaveSession stores WebAuthn session data with the given key.
func (db *SessionDB) SaveSession(key string, session *webauthn.SessionData) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.sessions[key] = session
}

// GetSession retrieves WebAuthn session data by its key.
func (db *SessionDB) GetSession(key string) (*webauthn.SessionData, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	session, ok := db.sessions[key]
	if !ok {
		return nil, nil
	}

	return session, nil
}

// DeleteSession removes WebAuthn session data with the given key.
func (db *SessionDB) DeleteSession(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.sessions, key)
}
