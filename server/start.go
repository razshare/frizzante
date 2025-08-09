package server

import (
	"context"
	"errors"
	"fmt"
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
							message := fmt.Sprintf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", r.Pattern, tag, g.Name)
							if conf.InfoLog != nil {
								conf.InfoLog.Println(message)
							} else {
								fmt.Println(message)
							}
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

		message := fmt.Sprintf("server bound to address %s; visit your application at http://%s\n", conf.Http.Addr, readableAddress)

		if conf.InfoLog != nil {
			conf.InfoLog.Printf(message)
		} else {
			fmt.Printf(message)
		}

		if exit {
			message = "cancelling server startup"
			if conf.InfoLog != nil {
				conf.InfoLog.Println(message)
			} else {
				fmt.Println(message)
			}
			return
		}

		serveError := http.ListenAndServe(conf.Http.Addr, conf.Http.Handler)
		if serveError != nil {
			if errors.Is(serveError, http.ErrServerClosed) {
				message = "shutting down server"
				if conf.InfoLog != nil {
					conf.InfoLog.Println(message)
				} else {
					fmt.Println(message)
				}
				return
			}
			conf.ErrorLog.Println(serveError)
		}
	}()

	go func() {
		if "" != conf.Certificate && "" != conf.Key {
			readableAddress := strings.Replace(conf.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)

			message := fmt.Sprintf("server bound to address %s; visit your application at https://%s", conf.Http.Addr, readableAddress)

			if conf.InfoLog != nil {
				conf.InfoLog.Println(message)
			} else {
				fmt.Println(message)
			}

			if exit {
				message = "cancelling server startup"
				if conf.InfoLog != nil {
					conf.InfoLog.Println(message)
				} else {
					fmt.Println(message)
				}
				return
			}

			serveError := http.ListenAndServeTLS(conf.SecureAddr, conf.Certificate, conf.Key, conf.Http.Handler)
			if serveError != nil {
				if errors.Is(serveError, http.ErrServerClosed) {
					message = "shutting down server"
					if conf.InfoLog != nil {
						conf.InfoLog.Println(message)
					} else {
						fmt.Println(message)
					}
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
