package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/razshare/frizzante/globals"
)

// New creates a new server.
func New() (server *Server) {
	return &Server{
		InfoLog:    log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		SecureAddr: "0.0.0.0:8383",
		PublicRoot: "app/dist/client",
		Channels: Channels{
			Stop: make(chan any, 1),
		},
		Server: &http.Server{
			Addr:           "0.0.0.0:8080",
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 3 * globals.MB,
			ErrorLog:       log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		},
	}
}
