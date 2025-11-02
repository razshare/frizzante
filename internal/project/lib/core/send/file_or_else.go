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

// FileOrElse sends the file requested by the client, or else falls back.
// Deprecated: use send.RequestedFile() instead.
func FileOrElse(client *clients.Client, orElse func()) {
	if client.WebSocket != nil {
		client.Options.ErrorLog.Println("file_or_else does not support web sockets", stack.Trace())
		return
	}

	if client.EventName != "" {
		client.Options.ErrorLog.Println("file_or_else does not support server sent events", stack.Trace())
		return
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

		http.ServeFile(client.Writer, &client.Request, name)
		return
	}

	if embeds.IsFile(client.Options.Efs, name) {
		var file fs.File
		var err error
		if file, err = client.Options.Efs.Open(name); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return
		}

		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return
		}

		if client.Writer.Header().Get("Content-Type") == "" {
			Header(client, "Content-Type", mime.Parse(name))
		}

		if client.Writer.Header().Get("Content-Length") == "" {
			Header(client, "Content-Length", fmt.Sprintf("%d", info.Size()))
		}

		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return
		}

		http.ServeContent(client.Writer, &client.Request, name, info.ModTime(), bytes.NewReader(buf))
		return
	}

	orElse()
}
