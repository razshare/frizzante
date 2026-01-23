package csr

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/servers"
)

// NewServer creates a new server.
func NewServer() (server *servers.Server) {
	return &servers.Server{
		InfoLog:        log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		SecureAddr:     "0.0.0.0:8383",
		Cors:           http.NewCrossOriginProtection(),
		Addr:           "0.0.0.0:8080",
		Handler:        http.NewServeMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2097152, // 2MB
		ErrorLog:       log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		Render:         NewFunction(),
		Channels: servers.Channels{
			Started: make(chan struct{}, 1),
			Ended:   make(chan struct{}, 1),
		},
	}
}
