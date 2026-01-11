package servers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/values"
)

// Snapshot generates static pages from a server.
func Snapshot(server *Server) {
	var err error
	go Start(server)
	<-server.Channels.Start
	server.Channels.Start <- values.None

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
		if response, err = http.Get(url); err != nil {
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
	}

	server.Channels.End <- values.None

	if err = embeds.CopyDirectory(server.Efs, "app/dist/client/assets", filepath.Join(".gen", "snapshot", "assets")); err != nil {
		return
	}

	return
}
