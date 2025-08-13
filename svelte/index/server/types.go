package server

import (
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/index"
)

type Server struct {
	*server.Server
	Index *index.Index
}
