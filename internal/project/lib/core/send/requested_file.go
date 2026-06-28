//go:build !use_disk

package send

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/embeds"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/mime"
)

// RequestedFile sends the file requested by the http.
//
// Returns false if connection is web sockets, server sent events
// or the file was not found.
func RequestedFile(writer http.ResponseWriter, request *http.Request, efs embed.FS, root string) (found bool, err error) {
	uri := strings.TrimPrefix(strings.ReplaceAll(filepath.Join(root, request.RequestURI), "/", string(filepath.Separator)), "/")
	embeddedFileName := strings.Join([]string{"app", "dist", "client", uri}, "/")
	if embeds.IsFile(efs, embeddedFileName) {
		var file fs.File
		if file, err = efs.Open(embeddedFileName); err != nil {
			return
		}
		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			return
		}
		header := writer.Header()
		if header.Get("Content-Type") == "" {
			header.Set("Content-Type", mime.Parse(embeddedFileName))
		}
		if writer.Header().Get("Content-Length") == "" {
			header.Set("Content-Length", fmt.Sprintf("%d", info.Size()))
		}
		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			return
		}
		http.ServeContent(writer, request, embeddedFileName, info.ModTime(), bytes.NewReader(buf))
		found = true
		return
	}
	fileName := filepath.Join("app", "dist", "client", strings.ReplaceAll(uri, "/", string(filepath.Separator)))
	if files.IsFile(fileName) {
		header := writer.Header()
		if writer.Header().Get("Content-Type") == "" {
			header.Set("Content-Type", mime.Parse(fileName))
		}
		http.ServeFile(writer, request, fileName)
		found = true
		return
	}
	return
}
