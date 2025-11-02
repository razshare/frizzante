package routes

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/guards"
)

type Route struct {
	Pattern string
	Handler func(client *clients.Client)
	Guards  []guards.Guard
}
