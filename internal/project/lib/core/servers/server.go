package servers

import (
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
)

type Server struct {
	http.Server
	Routes      []routes.Route
	SecureAddr  string
	Certificate string
	Key         string
	InfoLog     *log.Logger
	Cors        *http.CrossOriginProtection
}
