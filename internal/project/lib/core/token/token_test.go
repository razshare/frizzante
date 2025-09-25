package token

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	userId := "testuser123"

	// Generate token
	token1, err := Generate(userId)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token1 == "" {
		t.Error("Generated token is empty")
	}

	// Generate another token for same user
	token2, err := Generate(userId)
	if err != nil {
		t.Fatalf("Failed to generate second token: %v", err)
	}

	// Tokens should be different (due to random component)
	if token1 == token2 {
		t.Error("Two tokens generated for same user are identical")
	}
}

func TestTokenStore(t *testing.T) {
	// Create a new store for testing
	store := &TokenStore{
		tokens: make(map[string]TokenData),
	}

	userId := "testuser456"
	token := "test-token-123"
	ttl := 1 * time.Hour

	// Save token
	store.Save(token, userId, ttl)

	// Validate token
	data, valid := store.Validate(token)
	if !valid {
		t.Error("Token should be valid")
	}

	if data.UserId != userId {
		t.Errorf("Expected UserId %s, got %s", userId, data.UserId)
	}

	// Test invalid token
	_, valid = store.Validate("invalid-token")
	if valid {
		t.Error("Invalid token should not be valid")
	}

	store.Remove(token)

	// Token should no longer be valid
	_, valid = store.Validate(token)
	if valid {
		t.Error("Removed token should not be valid")
	}
}

func TestExpiredToken(t *testing.T) {
	// Create a new store for testing
	store := &TokenStore{
		tokens: make(map[string]TokenData),
	}

	userId := "testuser789"
	token := "test-token-expired"

	// Save token with very short TTL
	store.Save(token, userId, 1*time.Millisecond)

	// Wait for token to expire
	time.Sleep(2 * time.Millisecond)

	// Token should be expired
	_, valid := store.Validate(token)
	if valid {
		t.Error("Expired token should not be valid")
	}
}

func TestCleanExpired(t *testing.T) {
	// Create a new store for testing
	store := &TokenStore{
		tokens: make(map[string]TokenData),
	}

	// Add some tokens with different expiration times
	store.Save("token1", "user1", 1*time.Millisecond)
	store.Save("token2", "user2", 1*time.Hour)
	store.Save("token3", "user3", 1*time.Millisecond)

	// Wait for some tokens to expire
	time.Sleep(2 * time.Millisecond)

	// Clean expired tokens
	store.CleanExpired()

	// Only token2 should remain
	_, valid1 := store.Validate("token1")
	_, valid2 := store.Validate("token2")
	_, valid3 := store.Validate("token3")

	if valid1 {
		t.Error("Expired token1 should have been cleaned")
	}

	if !valid2 {
		t.Error("Valid token2 should still exist")
	}

	if valid3 {
		t.Error("Expired token3 should have been cleaned")
	}
}
