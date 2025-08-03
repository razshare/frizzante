package servers

import (
	"embed"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/views"
	"log"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	Efs                 embed.FS
	Guards              []guards.Guard
	Routes              []routes.Route
	SessionArchive      archives.Archive
	Connections         map[string]*net.Conn
	InfoLog             *log.Logger
	ViewContainers      []*views.Container
	ViewContainersMutex *sync.Mutex
	SecureAddr          string
	PublicRoot          string
	Dotenv              string
	Certificate         string
	Key                 string
	http.Server
}
