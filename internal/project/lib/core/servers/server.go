package servers

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

type Server struct {
	http.Server
	Routes      []routes.Route
	SecureAddr  string
	Certificate string
	Key         string
	InfoLog     *log.Logger
	Cors        *http.CrossOriginProtection
	Efs         embed.FS
	Render      renders.Render
}
