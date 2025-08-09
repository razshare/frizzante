package server

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/container"
	"net/http"
	"slices"
	"strings"
)

// Start starts a server from a configuration.
func Start(conf *Config) {
	cont := container.New(conf.Container)
	container.Start(cont)
	defer func() { go func() { cont.Channels.Stop <- 0 }() }()

	mux := conf.Http.Handler.(*http.ServeMux)

	if conf.Routes != nil {
		for _, r := range conf.Routes {
			mux.HandleFunc(r.Pattern, func(wrt http.ResponseWriter, req *http.Request) {
				con := &client.Client{
					Writer:  wrt,
					Request: req,
					Scope: client.Scope{
						Container: cont,
						EventId:   1,
						Status:    200,
					},
				}

				for _, tag := range r.Tags {
					for _, g := range conf.Guards {
						if !slices.Contains(g.Tags, tag) {
							continue
						}
						allow := false
						g.Handler(con, func() { allow = true })
						if !allow {
							conf.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", r.Pattern, tag, g.Name)
							return
						}
					}
				}

				r.Handler(con)
			})
		}
	}

	var exit bool

	go func() {
		readableAddress := strings.Replace(conf.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		conf.InfoLog.Printf("server bound to address %s; visit your application at http://%s", conf.Http.Addr, readableAddress)
		if exit {
			conf.InfoLog.Println("cancelling server startup")
			return
		}
		serveError := http.ListenAndServe(conf.Http.Addr, conf.Http.Handler)
		if serveError != nil {
			if errors.Is(serveError, http.ErrServerClosed) {
				conf.InfoLog.Println("shutting down server")
				return
			}
			conf.ErrorLog.Println(serveError)
		}
	}()

	go func() {
		if "" != conf.Certificate && "" != conf.Key {
			readableAddress := strings.Replace(conf.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			conf.InfoLog.Printf("server bound to address %s; visit your application at https://%s", conf.Http.Addr, readableAddress)
			if exit {
				conf.InfoLog.Println("cancelling server startup")
				return
			}
			serveError := http.ListenAndServeTLS(conf.SecureAddr, conf.Certificate, conf.Key, conf.Http.Handler)
			if serveError != nil {
				if errors.Is(serveError, http.ErrServerClosed) {
					conf.InfoLog.Println("shutting down server")
					return
				}
				conf.ErrorLog.Println(serveError)
			}
		}
	}()

	<-conf.Channels.Stop
	exit = true

	if err := conf.Http.Shutdown(context.Background()); err != nil {
		conf.ErrorLog.Println(err)
	}
}
