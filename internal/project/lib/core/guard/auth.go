package guard

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/tag"
	"github.com/razshare/frizzante/internal/project/lib/core/token"
)

// Auth is a guard that validates authentication tokens
var Auth = Guard{
	Name: "Auth",
	Handler: func(c *client.Client, allow func()) {
		// Extract token from request
		authToken := c.ExtractToken()
		if authToken == "" {
			c.Status = 401
			send.Json(c, map[string]string{
				"error": "No authentication token provided",
			})
			return
		}

		// Validate token
		tokenData, valid := token.Store.Validate(authToken)
		if !valid {
			c.Status = 401
			send.Json(c, map[string]string{
				"error": "Invalid or expired token",
			})
			return
		}

		c.UserId = tokenData.UserId
		allow()
	},
	Tags: []tag.Tag{},
}

// OptionalAuth validates token if present but doesn't require it
var OptionalAuth = Guard{
	Name: "OptionalAuth",
	Handler: func(c *client.Client, allow func()) {
		// Extract token from request
		authToken := c.ExtractToken()
		if authToken != "" {
			// Validate token
			tokenData, valid := token.Store.Validate(authToken)
			if valid {
				c.UserId = tokenData.UserId
			}
		}

		// Always allow to proceed for optional auth
		allow()
	},
	Tags: []tag.Tag{},
}
