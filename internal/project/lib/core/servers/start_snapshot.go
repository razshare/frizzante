//go:build snapshot_servers

package servers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"main/lib/core/clients"
	"main/lib/core/views/render"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"main/lib/core/embeds"
	"main/lib/core/stack"
	"main/lib/core/values"
)

// Start generates static pages from a server.
func Start(server *Server) {
	var err error
	snapshot := func() {
		client := http.Client{}
		if err = os.RemoveAll(filepath.Join(".gen", "snapshot")); err != nil {
			server.ErrorLog.Println(err, stack.Trace())
			os.Exit(1)
			return
		}
		for _, route := range server.Routes {
			var path string
			var method string
			var parts []string
			if parts = strings.SplitN(route.Pattern, " ", 2); len(parts) < 2 {
				err = fmt.Errorf("could not generate snapshot; pattern must be composed of a verb and path separated by a blank space; received %s", route.Pattern)
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			path = parts[1]
			method = parts[0]
			if method != "GET" {
				err = fmt.Errorf("could not generate snapshot; only GET verbs are allowed; received %s", route.Pattern)
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			if !strings.HasPrefix(path, "/") {
				err = fmt.Errorf("snapshot path must be absolute and thus start with /, received %s", path)
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			address := strings.Replace(server.Addr, "0.0.0.0", "127.0.0.1", 1)
			address = strings.Replace(address, "::", "127.0.0.1", 1)
			if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
				address = fmt.Sprintf("http://%s", address)
			}
			url := fmt.Sprintf("%s%s", address, path)
			var response *http.Response
			if response, err = client.Get(url); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			var data []byte
			if data, err = io.ReadAll(response.Body); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			if err = os.MkdirAll(filepath.Join(".gen", "snapshot", path), os.ModePerm); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			fsPath := strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/", string(filepath.Separator))
			if err = os.WriteFile(filepath.Join(".gen", "snapshot", fsPath, "index.html"), data, os.ModePerm); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			var request *http.Request
			if request, err = http.NewRequest("GET", url, nil); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			request.Header.Add("Accept", "application/json")
			if response, err = client.Do(request); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			if data, err = io.ReadAll(response.Body); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
			if err = os.WriteFile(filepath.Join(".gen", "snapshot", fsPath, "data.json"), data, os.ModePerm); err != nil {
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
		}
		if err = embeds.CopyDirectory(server.Efs, "app/dist/client/assets", filepath.Join(".gen", "snapshot", "assets")); err != nil {
			return
		}
		if server.Channels.End != nil {
			server.Channels.End <- values.None
		}
	}
	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.Handler,
		ReadTimeout:    server.ReadTimeout,
		WriteTimeout:   server.WriteTimeout,
		MaxHeaderBytes: server.MaxHeaderBytes,
		ErrorLog:       server.ErrorLog,
	}
	var render_ render.Render
	if render_, err = render.New(); err != nil {
		httpServer.ErrorLog.Println(err)
		return
	}
	handler := server.Handler.(*http.ServeMux)
	for _, route := range server.Routes {
		handler.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			if err = server.Cors.Check(request); err != nil {
				server.ErrorLog.Println(err)
				return
			}
			client := &clients.Client{
				Writer:  writer,
				Request: *request,
				Options: clients.Options{
					ErrorLog: server.ErrorLog,
					InfoLog:  server.InfoLog,
					Efs:      server.Efs,
					Render:   render_,
				},
				EventId: 1,
				Status:  200,
			}
			for _, guard := range route.Guards {
				allow := false
				guard.Handler(client, func() { allow = true })
				if !allow {
					if guard.Name == "" {
						server.InfoLog.Printf("an unnamed guard blocked the request on route %s", route.Pattern)
					} else {
						server.InfoLog.Printf("guard %s blocked the request on route %s", guard.Name, route.Pattern)
					}
					return
				}
			}
			route.Handler(client)
			if client.Channels.End != nil {
				client.Channels.End <- values.None
			}
		})
	}
	address := strings.Replace(httpServer.Addr, "0.0.0.0:", "127.0.0.1:", 1)
	server.InfoLog.Printf("server bound to address %s; visit your application at http://%s", httpServer.Addr, address)
	go func() {
		go snapshot()
		if err = httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				server.InfoLog.Println("shutting down server")
				return
			}
			server.ErrorLog.Println(err, stack.Trace())
			os.Exit(1)
		}
	}()
	<-server.Channels.End
	if err = httpServer.Shutdown(context.Background()); err != nil {
		httpServer.ErrorLog.Println(err)
		os.Exit(1)
	}
	return
}
