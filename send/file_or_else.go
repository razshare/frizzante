package send

import (
	"bytes"
	"fmt"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/mime"
	"github.com/razshare/frizzante/stack"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileOrElse sends the file requested by the client, or else falls back.
func FileOrElse(c *client.Client, or func()) {
	if c.WebSocket != nil {
		c.Config.ErrorLog.Println("FileOrElse does not support web sockets")
		return
	}

	if c.EventName != "" {
		c.Config.ErrorLog.Println("FileOrElse does not support server sent events")
		return
	}

	var n string

	if strings.HasPrefix(c.Request.RequestURI, "/") {
		n = filepath.Join(c.Config.PublicRoot, c.Request.RequestURI[1:])
	} else {
		n = filepath.Join(c.Config.PublicRoot, c.Request.RequestURI)
	}

	if embeds.IsFile(c.Config.Efs, n) {
		var f fs.File
		var err error
		if f, err = c.Config.Efs.Open(n); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		var i os.FileInfo
		if i, err = f.Stat(); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		if "" == c.Writer.Header().Get("Content-Type") {
			Header(c, "Content-Type", mime.Parse(n))
		}

		if "" == c.Writer.Header().Get("Content-Length") {
			Header(c, "Content-Length", fmt.Sprintf("%d", i.Size()))
		}

		buf := make([]byte, i.Size())
		if _, err = f.Read(buf); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		http.ServeContent(c.Writer, c.Request, n, i.ModTime(), bytes.NewReader(buf))
		return
	}

	if files.IsFile(n) {
		if "" == c.Writer.Header().Get("Content-Type") {
			Header(c, "Content-Type", mime.Parse(n))
		}

		http.ServeFile(c.Writer, c.Request, n)
		return
	}

	or()
}
