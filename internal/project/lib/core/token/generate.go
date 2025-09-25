package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/razshare/frizzante/text"
)

// Generate creates a new authentication token
func Generate(userId string) (token string, err error) {
	// Generate random bytes for token uniqueness
	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Combine user ID, timestamp, and random data
	tokenData := fmt.Sprintf("%s:%d:%s",
		userId,
		time.Now().Unix(),
		base64.StdEncoding.EncodeToString(randomBytes))

	// Create SHA1 hash of the token data
	token = text.Sha1(tokenData)
	return token, nil
}
