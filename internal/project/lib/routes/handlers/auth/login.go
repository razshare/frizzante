package auth

import (
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/token"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *client.Client) {
	var req LoginRequest
	if !receive.Json(c, &req) {
		c.Status = 400
		send.Json(c, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// TODO: Validate credentials against database
	// This is a simplified example - in production, verify against stored credentials
	if req.Username == "" || req.Password == "" {
		c.Status = 400
		send.Json(c, map[string]string{
			"error": "Username and password are required",
		})
		return
	}

	// Generate authentication token
	authToken, err := token.Generate(req.Username)
	if err != nil {
		c.Status = 500
		send.Json(c, map[string]string{
			"error": "Failed to generate authentication token",
		})
		return
	}

	token.Store.Save(authToken, req.Username, 24*time.Hour)
	c.SetToken(authToken)

	c.Status = 200
	send.Json(c, map[string]interface{}{
		"success": true,
		"token":   authToken,
		"user":    req.Username,
		"message": "Login successful",
	})
}
