package client

import (
	"strings"
)

// ExtractToken extracts the authentication token from the request
// Checks in order: Authorization header, Cookie, Query parameter
func (c *Client) ExtractToken() string {
	if c.Token != "" {
		return c.Token
	}

	// Check Authorization header (Bearer token)
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			c.Token = parts[1]
			return c.Token
		}
	}

	// Check cookie
	if cookie, err := c.Request.Cookie("auth_token"); err == nil {
		c.Token = cookie.Value
		return c.Token
	}

	// Check query parameter
	if token := c.Request.URL.Query().Get("token"); token != "" {
		c.Token = token
		return c.Token
	}

	return ""
}

// SetToken sets the authentication token in the client and as a cookie
func (c *Client) SetToken(token string) {
	c.Token = token
	// Set as HttpOnly cookie for security
	c.Writer.Header().Set("Set-Cookie",
		"auth_token="+token+"; HttpOnly; Path=/; SameSite=Strict")
}

// ClearToken removes the authentication token
func (c *Client) ClearToken() {
	c.Token = ""
	c.UserId = ""
	// Clear the cookie
	c.Writer.Header().Set("Set-Cookie",
		"auth_token=; HttpOnly; Path=/; Max-Age=0; SameSite=Strict")
}
