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
func Start(s *Server) {
	mux := s.Http.Handler.(*http.ServeMux)
	for _, r := range s.Routes {
		mux.HandleFunc(r.Pattern, func(wrt http.ResponseWriter, req *http.Request) {
			con := &client.Client{
				Writer:  wrt,
				Request: req,
				Scope: client.Scope{
					Render:     s.Render,
					ErrorLog:   s.ErrorLog,
					InfoLog:    s.InfoLog,
					PublicRoot: s.PublicRoot,
					Efs:        s.Efs,
					EventId:    1,
					Status:     200,
				},
			}

			for _, tag := range r.Tags {
				for _, g := range s.Guards {
					if !slices.Contains(g.Tags, tag) {
						continue
					}
					allow := false
					g.Handler(con, func() { allow = true })
					if !allow {
						s.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", r.Pattern, tag, g.Name)
						return
					}
				}
			}

			r.Handler(con)
		})
	}

	var exit bool

	go func() {
		haddr := strings.Replace(s.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		s.InfoLog.Printf("server bound to address %s; visit your application at http://%s", s.Http.Addr, haddr)
		if exit {
			s.InfoLog.Println("cancelling server startup")
			return
		}
		err := http.ListenAndServe(s.Http.Addr, s.Http.Handler)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				s.InfoLog.Println("shutting down server")
				return
			}
			s.ErrorLog.Println(err)
		}
	}()

	go func() {
		if "" != s.Certificate && "" != s.Key {
			haddr := strings.Replace(s.Http.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			s.InfoLog.Printf("server bound to address %s; visit your application at https://%s", s.Http.Addr, haddr)
			if exit {
				s.InfoLog.Println("cancelling server startup")
				return
			}
			err := http.ListenAndServeTLS(s.SecureAddr, s.Certificate, s.Key, s.Http.Handler)
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					s.InfoLog.Println("shutting down server")
					return
				}
				s.ErrorLog.Println(err)
			}
		}
	}()

	<-s.Channels.Stop
	exit = true

	if err := s.Http.Shutdown(context.Background()); err != nil {
		s.ErrorLog.Println(err)
	}
}
