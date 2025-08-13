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
	Render      func(v view.View) (string, error)
	Guards      []guard.Guard
	Routes      []route.Route
	Http        *http.Server
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	Channels    Channels
	Efs         embed.FS
	PublicRoot  string
	SecureAddr  string
	Certificate string
	Key         string
}

type Channels struct {
	Stop chan any
}
