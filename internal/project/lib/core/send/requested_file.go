//go:build !use_disk

package send

import (
	"bytes"
	"fmt"
	"io/fs"
	http_ "net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
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
		logs.Errorf(
			http,
			"send.RequestedFile: web sockets are not supported\n%s",
			stack.Trace(),
		)
		return false
	}
	if http.EventName != "" {
		logs.Errorf(
			http,
			"send.RequestedFile: server sent events are not supported\n%s",
			stack.Trace(),
		)
		return false
	}
	uri := http.Request.RequestURI
	if strings.HasPrefix(uri, "/") {
		uri = uri[1:]
	}
	embeddedFileName := strings.Join([]string{"app", "dist", "client", uri}, "/")
	if embeds.IsFile(http.Efs, embeddedFileName) {
		var file fs.File
		var err error
		if file, err = http.Efs.Open(embeddedFileName); err != nil {
			logs.Errorf(
				http,
				"send.RequestedFile: failed to open embedded file: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			logs.Errorf(
				http,
				"send.RequestedFile: failed to stat embedded file: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		if http.Writer.Header().Get("Content-Type") == "" {
			Header(http, "Content-Type", mime.Parse(embeddedFileName))
		}
		if http.Writer.Header().Get("Content-Length") == "" {
			Header(http, "Content-Length", fmt.Sprintf("%d", info.Size()))
		}
		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			logs.Errorf(
				http,
				"send.RequestedFile: failed to read embedded file: %v\n%s",
				err,
				stack.Trace(),
			)
			return false
		}
		http_.ServeContent(http.Writer, &http.Request, embeddedFileName, info.ModTime(), bytes.NewReader(buf))
		return true
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
