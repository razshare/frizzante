package routes

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

type Handler = func(client *clients.Client)
