package server

import (
	"embed"
	"github.com/razshare/frizzante/guard"
	"github.com/razshare/frizzante/route"
	"github.com/razshare/frizzante/view"
	"log"
	"net/http"
)

type Server struct {
	*http.Server
	PublicRoot  string
	SecureAddr  string
	Render      func(v view.View) (string, error)
	Guards      []guard.Guard
	Routes      []route.Route
	InfoLog     *log.Logger
	Channels    Channels
	Efs         embed.FS
	Certificate string
	Key         string
}

type Channels struct {
	Stop chan any
}
