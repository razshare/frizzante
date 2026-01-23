package servers

import (
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/views/render"
)

type Server struct {
	Routes         []routes.Route
	Channels       Channels
	Addr           string
	SecureAddr     string
	Certificate    string
	Key            string
	Handler        http.Handler
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	Cors           *http.CrossOriginProtection
	Efs            embed.FS
	Render         render.Function
}
