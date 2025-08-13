package server

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/client"
	"net/http"
	"slices"
	"strings"
)

// Start starts a server from a configuration.
func Start(c *Config) {
	mux := c.Http.Handler.(*http.ServeMux)
	for _, r := range c.Routes {
		mux.HandleFunc(r.Pattern, func(wrt http.ResponseWriter, req *http.Request) {
			con := &client.Client{
				Writer:  wrt,
				Request: req,
				Scope: client.Scope{
					Render:     c.Render,
					ErrorLog:   c.ErrorLog,
					InfoLog:    c.InfoLog,
					PublicRoot: c.PublicRoot,
					Efs:        c.Efs,
					EventId:    1,
					Status:     200,
				},
			}

			for _, tag := range r.Tags {
				for _, g := range c.Guards {
					if !slices.Contains(g.Tags, tag) {
						continue
					}
					allow := false
					g.Handler(con, func() { allow = true })
					if !allow {
						c.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", r.Pattern, tag, g.Name)
						return
					}
				}
			}

			r.Handler(con)
		})
	}

	var exit bool

	go func() {
		haddr := strings.Replace(c.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		c.InfoLog.Printf("server bound to address %s; visit your application at http://%s", c.Http.Addr, haddr)
		if exit {
			c.InfoLog.Println("cancelling server startup")
			return
		}
		err := http.ListenAndServe(c.Http.Addr, c.Http.Handler)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				c.InfoLog.Println("shutting down server")
				return
			}
			c.ErrorLog.Println(err)
		}
	}()

	go func() {
		if "" != c.Certificate && "" != c.Key {
			haddr := strings.Replace(c.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			c.InfoLog.Printf("server bound to address %s; visit your application at https://%s", c.Http.Addr, haddr)
			if exit {
				c.InfoLog.Println("cancelling server startup")
				return
			}
			err := http.ListenAndServeTLS(c.SecureAddr, c.Certificate, c.Key, c.Http.Handler)
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					c.InfoLog.Println("shutting down server")
					return
				}
				c.ErrorLog.Println(err)
			}
		}
	}()

	<-c.Channels.Stop
	exit = true

	if err := c.Http.Shutdown(context.Background()); err != nil {
		c.ErrorLog.Println(err)
	}
}
