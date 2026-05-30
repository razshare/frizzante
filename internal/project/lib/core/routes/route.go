package routes

import "github.com/razshare/frizzante/internal/project/lib/core/guards"

type Route struct {
	Pattern string
	Handler Handler
	Guards  []guards.Guard
}
