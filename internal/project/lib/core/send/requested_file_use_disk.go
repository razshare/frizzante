//go:build use_disk

package send

import (
	http_ "net/http"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/mime"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// RequestedFile sends the file requested by the http.
//
// Returns false if connection is web sockets, server sent events
// or the file was not found.
func RequestedFile(http *scopes.Http) bool {
	if http.WebSocket != nil {
		logs.Errorf(http, "send.RequestedFile() does not support web sockets\n%s", stack.Trace())
		return false
	}
	if http.EventName != "" {
		logs.Errorf(http, "send.RequestedFile() does not support server sent events\n%s", stack.Trace())
		return false
	}
	uri := http.Request.RequestURI
	if strings.HasPrefix(uri, "/") {
		uri = uri[1:]
	}
	fileName := filepath.Join("app", "dist", "client", strings.ReplaceAll(uri, "/", string(filepath.Separator)))
	if files.IsFile(fileName) {
		if http.Writer.Header().Get("Content-Type") == "" {
			Header(http, "Content-Type", mime.Parse(fileName))
		}
		http_.ServeFile(http.Writer, &http.Request, fileName)
		return true
	}
	return false
}
