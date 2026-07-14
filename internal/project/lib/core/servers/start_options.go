package servers

import (
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
)

type StartOptions struct {
	Address     string
	Certificate string
	Key         string
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
	Routes      []routes.Route
	Cors        *http.CrossOriginProtection
	BeforeStart func(server *http.Server)
	AfterStop   func(server *http.Server)
}
