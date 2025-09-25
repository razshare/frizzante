package routes

import (
	"github.com/razshare/frizzante/internal/project/lib/core/guard"
	"github.com/razshare/frizzante/internal/project/lib/core/route"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/auth"
)

// AuthRoutes returns all authentication-related routes
func AuthRoutes() []route.Route {
	return []route.Route{
		// Public routes (no authentication required)
		{
			Pattern: "POST /api/auth/login",
			Handler: auth.Login,
		},
		{
			Pattern: "POST /api/auth/logout",
			Handler: auth.Logout,
		},
		// Protected routes (authentication required)
		{
			Pattern: "GET /api/auth/profile",
			Handler: guard.With(auth.Profile, guard.Auth),
		},
	}
}
