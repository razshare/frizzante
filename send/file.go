package send

import (
	"embed"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/mime"
	"github.com/razshare/frizzante/stack"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// EmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func EmbeddedFileOrElse(c *client.Client, efs embed.FS, or func()) {
	fileName := c.Scope.PublicRoot + c.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(efs, fileName) || embeds.IsDirectory(efs, fileName) {
		or()
		return
	}

	reader, readerInfo, readerError := embeds.NewFileReader(efs, fileName)
	if readerError != nil {
		c.Scope.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.Scope.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Scope.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.Scope.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.Scope.ErrorLog.Println(writeError, stack.Trace())
			return
		}
		return
	}

	if "" != c.Scope.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Scope.ErrorLog.Println(readError, stack.Trace())
			return
		}

		EventContent(c, data)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", mime.Parse(fileName))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		Header(c, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(c.Writer, c.Request, fileName, readerInfo.ModTime(), reader)
}

// FileOrElse sends the file requested by the client, or else falls back.
func FileOrElse(c *client.Client, or func()) {
	fileName := filepath.Join(c.Scope.PublicRoot, c.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		EmbeddedFileOrElse(c, c.Scope.Efs, or)
		return
	}

	reader, readerInfo, readerError := files.NewFileReader(fileName)
	if readerError != nil {
		c.Scope.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.Scope.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Scope.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.Scope.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.Scope.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	if "" != c.Scope.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Scope.ErrorLog.Println(readError, stack.Trace())
			return
		}

		EventContent(c, data)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", mime.Parse(fileName))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		Header(c, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(c.Writer, c.Request, fileName, readerInfo.ModTime(), reader)
}
