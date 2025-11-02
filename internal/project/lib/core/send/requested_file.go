//go:build !dev

package send

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/mime"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// RequestedFile sends the file requested by the client.
//
// Returns false if connection is web sockets, server sent events
// or the file was not found.
func RequestedFile(client *clients.Client) bool {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("send.RequestedFile() does not support web sockets", stack.Trace())
		return false
	}

	if client.EventName != "" {
		client.Config.ErrorLog.Println("send.RequestedFile() does not support server sent events", stack.Trace())
		return false
	}

	var name string

	if strings.HasPrefix(client.Request.RequestURI, "/") {
		name = filepath.Join("app", "dist", "client", client.Request.RequestURI[1:])
	} else {
		name = filepath.Join("app", "dist", "client", client.Request.RequestURI)
	}

	if files.IsFile(name) {
		if client.Writer.Header().Get("Content-Type") == "" {
			Header(client, "Content-Type", mime.Parse(name))
		}

		http.ServeFile(client.Writer, client.Request, name)
		return true
	}

	if embeds.IsFile(client.Config.Efs, name) {
		var file fs.File
		var err error
		if file, err = client.Config.Efs.Open(name); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}

		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}

		if client.Writer.Header().Get("Content-Type") == "" {
			Header(client, "Content-Type", mime.Parse(name))
		}

		if client.Writer.Header().Get("Content-Length") == "" {
			Header(client, "Content-Length", fmt.Sprintf("%d", info.Size()))
		}

		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}

		http.ServeContent(client.Writer, client.Request, name, info.ModTime(), bytes.NewReader(buf))
		return true
	}

	return false
}
