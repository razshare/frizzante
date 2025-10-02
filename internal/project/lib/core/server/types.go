package server

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/guard"
	"github.com/razshare/frizzante/internal/project/lib/core/route"
)

type Server struct {
	*http.Server
	Guards      []guard.Guard
	Routes      []route.Route
	App         string
	PublicRoot  string
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
