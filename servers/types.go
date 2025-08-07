package servers

import (
	"embed"
	"github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
	"net/http"
)

type Server struct {
	http.Server
	AppConfiguration apps.Configuration
	SessionArchive   archives.Archive
	Guards           []guards.Guard
	Routes           []routes.Route
	InfoLog          *log.Logger
	Efs              embed.FS
	Stop             chan any
	SecureAddr       string
	Certificate      string
	PublicRoot       string
	Dotenv           string
	Key              string
}
