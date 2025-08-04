package servers

import (
	"embed"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/containers"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
	"net"
	"net/http"
)

type Server struct {
	http.Server
	SecureAddr     string
	PublicRoot     string
	Dotenv         string
	Certificate    string
	Key            string
	Efs            embed.FS
	InfoLog        *log.Logger
	Guards         []guards.Guard
	Routes         []routes.Route
	SessionArchive archives.Archive
	Connections    map[string]*net.Conn
	ViewContainer  *containers.ViewContainer
}
