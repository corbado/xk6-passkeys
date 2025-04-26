package main

import (
	"sync"
)

// UserDB provides thread-safe in-memory storage for User objects.
type UserDB struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewUserDB creates a new instance of UserDB with an empty user map.
func NewUserDB() *UserDB {
	return &UserDB{
		users: make(map[string]*User),
	}
}

// GetUser retrieves a user by their username from the database.
func (db *UserDB) GetUser(name string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.users[name]
	if !ok {
		return nil, nil
	}

	return user, nil
}

// PutUser stores a user in the database using their username as the key.
func (db *UserDB) PutUser(user *User) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.users[user.name] = user
}
