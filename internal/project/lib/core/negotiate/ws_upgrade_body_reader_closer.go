package negotiate

import (
	"io"

	"github.com/gorilla/websocket"
)

type WsUpgradeBodyReaderCloser struct {
	io.ReadCloser
	WebSocketConnection *websocket.Conn
	Reader              io.Reader
}

func (readerCloser WsUpgradeBodyReaderCloser) Read(buffer []byte) (int, error) {
	_, reader, _ := readerCloser.WebSocketConnection.NextReader()
	return reader.Read(buffer)
}

func (readerCloser WsUpgradeBodyReaderCloser) Close() error {
	return readerCloser.Close()
}
