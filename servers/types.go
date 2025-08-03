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
	Efs            embed.FS
	Guards         []guards.Guard
	Routes         []routes.Route
	SessionArchive archives.Archive
	Connections    map[string]*net.Conn
	InfoLog        *log.Logger
	AppRoot        string
	ServerJs       string
	IndexHtml      string
	SecureAddr     string
	PublicRoot     string
	Dotenv         string
	Certificate    string
	Key            string
	http.Server
}
