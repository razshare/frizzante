package auth

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

// Profile returns the authenticated user's profile
// This handler should be protected by the Auth guard
func Profile(c *client.Client) {
	// UserId is populated by the Auth guard
	if c.UserId == "" {
		c.Status = 401
		send.Json(c, map[string]string{
			"error": "Not authenticated",
		})
		return
	}

	// TODO: Fetch user profile from database
	// This is a simplified example
	c.Status = 200
	send.Json(c, map[string]interface{}{
		"user_id": c.UserId,
		"message": "Profile retrieved successfully",
		// Add more profile fields as needed
	})
}
