package main

import (
	"sync"
)

// UserDB represents an in-memory user database
type UserDB struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewUserDB creates a new UserDB instance
func NewUserDB() *UserDB {
	return &UserDB{
		users: make(map[string]*User),
	}
}

// GetUser returns a *User by the user's username
func (db *UserDB) GetUser(name string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	user, ok := db.users[name]
	if !ok {
		return nil, nil
	}
	return user, nil
}

// PutUser stores a new user by the user's username
func (db *UserDB) PutUser(user *User) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.users[user.name] = user
}
