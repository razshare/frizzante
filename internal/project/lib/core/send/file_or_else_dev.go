//go:build dev

package send

import (
	"net/http"
	"path/filepath"
	"strings"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/mime"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FileOrElse sends the file requested by the client, or else falls back.
func FileOrElse(client *_client.Client, orElse func()) {
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

	if files.IsFile(name) {
		if "" == client.Writer.Header().Get("Content-Type") {
			Header(client, "Content-Type", mime.Parse(name))
		}

		http.ServeFile(client.Writer, client.Request, name)
		return
	}

	orElse()
}
