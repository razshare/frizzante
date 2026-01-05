package servers

import (
	"log"
	"net/http"
	"os"
	"time"
)

// New creates a new server.
func New() (server *Server) {
	return &Server{
		InfoLog:    log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		SecureAddr: "0.0.0.0:8383",
		Cors:       http.NewCrossOriginProtection(),
		Channels: Channels{
			Start: make(chan struct{}, 1),
			End:   make(chan struct{}, 1),
		},
		Addr:           "0.0.0.0:8080",
		Handler:        http.NewServeMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2097152, // 2MB
		ErrorLog:       log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
	}
}
