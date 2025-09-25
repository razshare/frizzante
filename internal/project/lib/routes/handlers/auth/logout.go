package auth

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/token"
)

func Logout(c *client.Client) {
	authToken := c.ExtractToken()

	if authToken != "" {
		token.Store.Remove(authToken)
	}

	c.ClearToken()

	c.Status = 200
	send.Json(c, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}
