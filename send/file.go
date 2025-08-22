package send

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/mime"
	"github.com/razshare/frizzante/stack"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileOrElse sends the file requested by the client, or else falls back.
func FileOrElse(client *client.Client, or func()) {
	var name = filepath.Join(client.Config.PublicRoot, client.Request.RequestURI)
	var reader *bytes.Reader
	var info os.FileInfo
	var err error

	if embeds.IsFile(client.Config.Efs, name) && !embeds.IsDirectory(client.Config.Efs, name) {
		reader, info, err = embeds.NewFileReader(client.Config.Efs, strings.ReplaceAll(name, "\\", "//"))
	} else if files.IsFile(name) && !files.IsDirectory(name) {
		reader, info, err = files.NewFileReader(name)
	} else {
		or()
		return
	}

	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	if client.WebSocket != nil {
		var data []byte
		if data, err = io.ReadAll(reader); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		if err = client.WebSocket.WriteMessage(websocket.TextMessage, data); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	}

	if "" != client.EventName {
		var data []byte
		if data, err = io.ReadAll(reader); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		EventContent(client, data)
		return
	}

	if "" == client.Writer.Header().Get("Content-Type") {
		Header(client, "Content-Type", mime.Parse(name))
	}

	if "" == client.Writer.Header().Get("Content-Length") {
		Header(client, "Content-Length", fmt.Sprintf("%d", info.Size()))
	}

	http.ServeContent(client.Writer, client.Request, name, info.ModTime(), reader)
}
