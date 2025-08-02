package servers

import (
	"embed"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Guards         []guards.Guard
	Routes         []routes.Route
	Efs            embed.FS
	SessionArchive archives.Archive
	Connections    map[string]*net.Conn
	InfoLog        *log.Logger
	SecureAddr     string
	PublicRoot     string
	AppRoot        string
	ServerJs       string
	IndexHtml      string
	Dotenv         string
	Certificate    string
	Key            string
	http.Server
}
