package server

import (
	"github.com/razshare/frizzante/globals"
	"log"
	"net/http"
	"os"
	"time"
)

// New creates a new server.
func New() (s *Server) {
	ilog := log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	elog := log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	return &Server{
		InfoLog:  ilog,
		ErrorLog: elog,
		Channels: Channels{
			Stop: make(chan any, 1),
		},
		Http: &Http{
			SecureAddr: "0.0.0.0:8383",
			PublicRoot: "app/dist/client",
			Server: &http.Server{
				Addr:           "0.0.0.0:8080",
				Handler:        http.NewServeMux(),
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 3 * globals.MB,
				ErrorLog:       elog,
			},
		},
	}
}
