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
		client.Options.ErrorLog.Println("send.RequestedFile() does not support web sockets", stack.Trace())
		return false
	}

	if client.EventName != "" {
		client.Options.ErrorLog.Println("send.RequestedFile() does not support server sent events", stack.Trace())
		return false
	}

	uri := client.Request.RequestURI
	if strings.HasPrefix(uri, "/") {
		uri = uri[1:]
	}

	embeddedFileName := strings.Join([]string{"app", "dist", "client", uri}, "/")

	if embeds.IsFile(client.Options.Efs, embeddedFileName) {
		var file fs.File
		var err error
		if file, err = client.Options.Efs.Open(embeddedFileName); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return false
		}

		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return false
		}

		if client.Writer.Header().Get("Content-Type") == "" {
			Header(client, "Content-Type", mime.Parse(embeddedFileName))
		}

		if client.Writer.Header().Get("Content-Length") == "" {
			Header(client, "Content-Length", fmt.Sprintf("%d", info.Size()))
		}

		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return false
		}

		http.ServeContent(client.Writer, &client.Request, embeddedFileName, info.ModTime(), bytes.NewReader(buf))
		return true
	}

	fileName := filepath.Join("app", "dist", "client", strings.ReplaceAll(uri, "/", string(filepath.Separator)))

	if files.IsFile(fileName) {
		if client.Writer.Header().Get("Content-Type") == "" {
			Header(client, "Content-Type", mime.Parse(fileName))
		}

		http.ServeFile(client.Writer, &client.Request, fileName)
		return true
	}

	return false
}
