package servers

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/views"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Guards         []guards.Guard
	Routes         []routes.Route
	SessionArchive archives.Archive
	ViewContainer  *views.Container
	Connections    map[string]*net.Conn
	InfoLog        *log.Logger
	SecureAddr     string
	PublicRoot     string
	Dotenv         string
	Certificate    string
	Key            string
	http.Server
}
