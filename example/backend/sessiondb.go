package main

import (
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
)

type SessionDB struct {
	sessions map[string]*webauthn.SessionData
	mu       sync.RWMutex
}

func NewSessionDB() *SessionDB {
	return &SessionDB{
		sessions: make(map[string]*webauthn.SessionData),
	}
}

func (db *SessionDB) SaveSession(key string, session *webauthn.SessionData) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.sessions[key] = session
}

func (db *SessionDB) GetSession(key string) (*webauthn.SessionData, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	session, ok := db.sessions[key]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (db *SessionDB) DeleteSession(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.sessions, key)
}
