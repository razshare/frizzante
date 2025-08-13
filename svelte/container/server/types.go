package server

import (
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/container"
)

type Server struct {
	*server.Server
	Container *container.Container
}
