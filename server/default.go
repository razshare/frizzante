package server

import (
	"github.com/razshare/frizzante/container"
	"github.com/razshare/frizzante/globals"
	"log"
	"net/http"
	"os"
	"time"
)

// Default creates a default server configuration.
func Default() *Config {
	infoLog := log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	return &Config{
		InfoLog:    infoLog,
		ErrorLog:   errorLog,
		SecureAddr: "0.0.0.0:8383",
		Channels: Channels{
			Stop: make(chan any, 1),
		},
		Http: &http.Server{
			Addr:           "0.0.0.0:8080",
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 3 * globals.MB,
			ErrorLog:       errorLog,
		},
		Container: container.Config{
			Root:        "app",
			Script:      "app/dist/server.js",
			Document:    "app/dist/client/index.html",
			PublicRoot:  "app/dist/client",
			Parallels:   2,
			InfoLog:     infoLog,
			ErrorLog:    errorLog,
			Development: os.Getenv("DEV") == "1",
		},
	}
}
