package server

import (
	"embed"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/index"
)

func Default(efs embed.FS) *Server {
	// New.
	s := server.New()
	i := index.New(efs)

	// Index.
	i.ErrorLog = s.ErrorLog

	// Server.
	s.Efs = efs
	s.Render = i.Render

	// Server.
	return &Server{
		Server: s,
		Index:  i,
	}
}
