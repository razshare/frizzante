package send

import (
	"bytes"
	"fmt"
	"io/fs"
	"main/lib/core/client"
	"main/lib/core/embeds"
	"main/lib/core/files"
	"main/lib/core/mime"
	"main/lib/core/stack"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileOrElse sends the file requested by the client, or else falls back.
func FileOrElse(client *client.Client, or func()) {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("file_or_else does not support web sockets", stack.Trace())
		return
	}

	if client.EventName != "" {
		client.Config.ErrorLog.Println("file_or_else does not support server sent events", stack.Trace())
		return
	}

	var name string

	if strings.HasPrefix(client.Request.RequestURI, "/") {
		name = filepath.Join(client.Config.PublicRoot, client.Request.RequestURI[1:])
	} else {
		name = filepath.Join(client.Config.PublicRoot, client.Request.RequestURI)
	}

	if embeds.IsFile(client.Config.Efs, name) {
		var file fs.File
		var err error
		if file, err = client.Config.Efs.Open(name); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		if "" == client.Writer.Header().Get("Content-Type") {
			Header(client, "Content-Type", mime.Parse(name))
		}

		if "" == client.Writer.Header().Get("Content-Length") {
			Header(client, "Content-Length", fmt.Sprintf("%d", info.Size()))
		}

		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		http.ServeContent(client.Writer, client.Request, name, info.ModTime(), bytes.NewReader(buf))
		return
	}

	if files.IsFile(name) {
		if "" == client.Writer.Header().Get("Content-Type") {
			Header(client, "Content-Type", mime.Parse(name))
		}

		http.ServeFile(client.Writer, client.Request, name)
		return
	}

	or()
}
