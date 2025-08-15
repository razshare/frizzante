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
func FileOrElse(c *client.Client, or func()) {
	var n = filepath.Join(c.Scope.PublicRoot, c.Request.RequestURI)
	var r *bytes.Reader
	var i os.FileInfo
	var err error

	if embeds.IsFile(c.Scope.Efs, n) && !embeds.IsDirectory(c.Scope.Efs, n) {
		r, i, err = embeds.NewFileReader(c.Scope.Efs, strings.ReplaceAll(n, "\\", "//"))
	} else if files.IsFile(n) && !files.IsDirectory(n) {
		r, i, err = files.NewFileReader(n)
	} else {
		or()
		return
	}

	if err != nil {
		c.Scope.ErrorLog.Println(err, stack.Trace())
		return
	}

	if c.Scope.WebSocket != nil {
		var d []byte
		d, err = io.ReadAll(r)
		if err != nil {
			c.Scope.ErrorLog.Println(err, stack.Trace())
			return
		}

		err = c.Scope.WebSocket.WriteMessage(websocket.TextMessage, d)
		if err != nil {
			c.Scope.ErrorLog.Println(err, stack.Trace())
			return
		}
	}

	if "" != c.Scope.EventName {
		var d []byte
		d, err = io.ReadAll(r)
		if err != nil {
			c.Scope.ErrorLog.Println(err, stack.Trace())
			return
		}

		EventContent(c, d)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", mime.Parse(n))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		Header(c, "Content-Length", fmt.Sprintf("%d", i.Size()))
	}

	http.ServeContent(c.Writer, c.Request, n, i.ModTime(), r)
}
