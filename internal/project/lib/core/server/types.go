package server

import (
	"embed"
	"log"
	"main/lib/core/guard"
	"main/lib/core/route"
	_view "main/lib/core/view"
	"net/http"
)

type Server struct {
	*http.Server
	PublicRoot  string
	SecureAddr  string
	Guards      []guard.Guard
	Routes      []route.Route
	InfoLog     *log.Logger
	Channels    Channels
	Efs         embed.FS
	Certificate string
	Key         string
	Render      func(view _view.View) (html string, err error)
}

type Channels struct {
	Stop chan any
}
