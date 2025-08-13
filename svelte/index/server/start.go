package server

import (
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/index"
)

func Start(a *Server) {
	go index.Start(a.Index)
	server.Start(a.Server)
}
