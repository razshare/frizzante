package send

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type SseUpgradeResponseWriter struct {
	http.ResponseWriter
	EventName string
	EventId   int
}

func (writer *SseUpgradeResponseWriter) Write(data []byte) (written int, err error) {
	var count int
	meta := fmt.Sprintf("id: %d\r\nevent: %s\r\n", writer.EventId, writer.EventName)
	if count, err = writer.ResponseWriter.Write([]byte(meta)); err != nil {
		return
	}
	written += count
	for _, line := range bytes.Split(data, []byte("\r\n")) {
		if count, err = writer.ResponseWriter.Write([]byte("data: ")); err != nil {
			return
		}
		written += count
		if count, err = writer.ResponseWriter.Write(line); err != nil {
			return
		}
		written += count
		if count, err = writer.ResponseWriter.Write([]byte("\r\n")); err != nil {
			return
		}
		written += count
	}
	if count, err = writer.ResponseWriter.Write([]byte("\r\n")); err != nil {
		return
	}
	written += count
	var ok bool
	var flusher http.Flusher
	if flusher, ok = writer.ResponseWriter.(http.Flusher); !ok {
		err = errors.New("send.EventContent: could not retrieve flusher")
		return
	}
	flusher.Flush()
	writer.EventId++
	return
}
