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
	Http        *Http
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	Channels    Channels
	Efs         embed.FS
	Certificate string
	Key         string
}

type Http struct {
	*http.Server
	PublicRoot string
	SecureAddr string
}

type Channels struct {
	Stop chan any
}
