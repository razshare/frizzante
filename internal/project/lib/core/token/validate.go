package token

import (
	"sync"
	"time"
)

type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]TokenData
}

type TokenData struct {
	UserId    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Global token store (in production, use a database or Redis)
var Store = &TokenStore{
	tokens: make(map[string]TokenData),
}

// Save stores a token with its metadata
func (s *TokenStore) Save(token string, userId string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.tokens[token] = TokenData{
		UserId:    userId,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}
}

// Validate checks if a token is valid and returns its data
func (s *TokenStore) Validate(token string) (data TokenData, valid bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.tokens[token]
	if !exists {
		return TokenData{}, false
	}

	// Check if token has expired
	if time.Now().After(data.ExpiresAt) {
		// Clean up expired token
		go s.Remove(token)
		return TokenData{}, false
	}

	return data, true
}

// Remove deletes a token from the store
func (s *TokenStore) Remove(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
}

// CleanExpired removes all expired tokens
func (s *TokenStore) CleanExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, data := range s.tokens {
		if now.After(data.ExpiresAt) {
			delete(s.tokens, token)
		}
	}
}
