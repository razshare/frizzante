package servers

import (
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/guards"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
)

type StartOptions struct {
	SecureAddr  string
	Certificate string
	Key         string
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
	Routes      []routes.Route
	Guards      []guards.Guard
	Cors        *http.CrossOriginProtection
}
