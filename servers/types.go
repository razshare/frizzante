package servers

import (
	"embed"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Http          *http.Server
	Guards        []guards.Guard
	Routes        []routes.Route
	Efs           embed.FS
	Connections   map[string]*net.Conn
	InfoLog       *log.Logger
	Address       string
	SecureAddress string
	PublicRoot    string
	AppRoot       string
	ServerJs      string
	IndexHtml     string
	Dotenv        string
	Certificate   string
	Key           string
}
