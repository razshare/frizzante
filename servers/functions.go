package servers

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

// New creates a new server.
func New() *Server {
	infoLog := log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	return &Server{
		InfoLog:    infoLog,
		SecureAddr: "0.0.0.0:8383",
		Stop:       make(chan any),
		Server: http.Server{
			Addr:           "0.0.0.0:8080",
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 3 * globals.MB,
			ErrorLog:       errorLog,
		},
		PublicRoot: "app/dist/client",
		AppConfig: apps.Config{
			Root:        "app",
			Script:      "app/dist/server.js",
			Document:    "app/dist/client/index.html",
			Parallels:   2,
			InfoLog:     infoLog,
			ErrorLog:    errorLog,
			Development: os.Getenv("DEV") == "1",
		},
	}
}

// Start starts the server.
func Start(s *Server) {
	application := apps.Start(s.AppConfig, s.Efs)
	defer func() { go func() { application.Stop <- 0 }() }()

	mux := s.Handler.(*http.ServeMux)

	for _, r := range s.Routes {
		mux.HandleFunc(r.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			connection := &connections.Connection{
				EventId:    1,
				Status:     200,
				Writer:     writer,
				Request:    request,
				App:        application,
				Efs:        s.Efs,
				ErrorLog:   s.ErrorLog,
				AppConfig:  &s.AppConfig,
				PublicRoot: s.PublicRoot,
			}

			for _, tag := range r.Tags {
				for _, guard := range s.Guards {
					if !slices.Contains(guard.Tags, tag) {
						continue
					}
					allow := false
					guard.Handler(connection, func() { allow = true })
					if !allow {
						s.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", r.Pattern, tag, guard.Name)
						return
					}
				}
			}

			r.Handler(connection)
		})
	}

	var exit bool

	go func() {
		readableAddress := strings.Replace(s.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		s.InfoLog.Printf("server bound to address %s; visit your application at http://%s", s.Addr, readableAddress)
		if exit {
			s.InfoLog.Printf("cancelling server startup")
			return
		}
		serveError := http.ListenAndServe(s.Addr, s.Handler)
		if serveError != nil {
			if errors.Is(serveError, http.ErrServerClosed) {
				s.InfoLog.Println("shutting down server")
				return
			}
			s.ErrorLog.Println(serveError)
		}
	}()

	go func() {
		if "" != s.Certificate && "" != s.Key {
			readableAddress := strings.Replace(s.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			s.InfoLog.Printf("server bound to address %s; visit your application at https://%s", s.Addr, readableAddress)
			if exit {
				s.InfoLog.Printf("cancelling server startup")
				return
			}
			serveError := http.ListenAndServeTLS(s.SecureAddr, s.Certificate, s.Key, s.Handler)
			if serveError != nil {
				if errors.Is(serveError, http.ErrServerClosed) {
					s.InfoLog.Printf("shutting down server")
					return
				}
				s.ErrorLog.Println(serveError)
			}
		}
	}()

	<-s.Stop
	exit = true

	if err := s.Shutdown(context.Background()); err != nil {
		s.ErrorLog.Println(err)
	}
}
