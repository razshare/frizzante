package server

import (
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/container"
)

func Start(a *Server) {
	go container.Start(a.Container)
	server.Start(a.Server)
}
