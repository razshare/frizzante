//go:build use_disk

package send

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/mime"
)

// RequestedFile sends the file requested by the http.
//
// Returns false if connection is web sockets, server sent events
// or the file was not found.
func RequestedFile(writer http.ResponseWriter, request *http.Request, _ embed.FS, root string) (found bool, err error) {
	uri := strings.ReplaceAll(filepath.Join(root, request.RequestURI), "/", string(filepath.Separator))
	if strings.HasPrefix(uri, "/") {
		uri = uri[1:]
	}
	fileName := filepath.Join("app", "dist", "client", uri)
	if files.IsFile(fileName) {
		found = true
		header := writer.Header()
		if header.Get("Content-Type") == "" {
			header.Set("Content-Type", mime.Parse(fileName))
		}
		http.ServeFile(writer, request, fileName)
		return
	}
	return
}
