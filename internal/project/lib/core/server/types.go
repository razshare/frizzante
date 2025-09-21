package server

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/guard"
	"github.com/razshare/frizzante/internal/project/lib/core/route"
	view "github.com/razshare/frizzante/internal/project/lib/core/view"
)

type Server struct {
	*http.Server
	Guards      []guard.Guard
	Routes      []route.Route
	PublicRoot  string
	SecureAddr  string
	Certificate string
	Key         string
	Channels    Channels
	InfoLog     *log.Logger
	Efs         embed.FS
	Render      func(view view.View) (html string, err error)
}

type Channels struct {
	Stop chan any
}
