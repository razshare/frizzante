package server

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/guard"
	"github.com/razshare/frizzante/internal/project/lib/core/route"
	_view "github.com/razshare/frizzante/internal/project/lib/core/view"
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
