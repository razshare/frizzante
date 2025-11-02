package servers

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
)

type Server struct {
	*http.Server
	Routes      []routes.Route
	SecureAddr  string
	Certificate string
	Key         string
	InfoLog     *log.Logger
	Cors        *http.CrossOriginProtection
	Channels    Channels
	Efs         embed.FS
}

type Channels struct {
	Stop chan any
}
