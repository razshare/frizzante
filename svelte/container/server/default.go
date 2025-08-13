package server

import (
	"embed"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/container"
)

func Default(efs embed.FS) *Server {
	// New.
	s := server.New()
	c := container.New(efs)

	// Container.
	c.ErrorLog = s.ErrorLog

	// Server.
	s.Efs = efs
	s.Render = c.Render

	// Server.
	return &Server{
		Server:    s,
		Container: c,
	}
}
