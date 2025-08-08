package send

import (
	"embed"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/conn"
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
func EmbeddedFileOrElse(c *conn.Conn, fs embed.FS, or func()) {
	fileName := c.Container.Config.PublicRoot + c.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(fs, fileName) || embeds.IsDirectory(fs, fileName) {
		or()
		return
	}

	reader, readerInfo, readerError := embeds.NewFileReader(fs, fileName)
	if readerError != nil {
		c.Container.Config.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}
		return
	}

	if "" != c.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
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
func FileOrElse(c *conn.Conn, or func()) {
	fileName := filepath.Join(c.Container.Config.PublicRoot, c.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		EmbeddedFileOrElse(c, c.Container.Config.Efs, or)
		return
	}

	reader, readerInfo, readerError := files.NewFileReader(fileName)
	if readerError != nil {
		c.Container.Config.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	if "" != c.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
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
