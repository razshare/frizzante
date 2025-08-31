package mock

import (
	"io"
	"main/lib/core/client"
	"main/lib/core/server"
	"net/http"
)

type ResponseWriter struct {
	MockHeader     http.Header
	MockStatusCode int
	MockBytes      []byte
}

func (m *ResponseWriter) Header() http.Header {
	return m.MockHeader
}

func (m *ResponseWriter) Write(bytes []byte) (int, error) {
	m.MockBytes = append(m.MockBytes, bytes...)
	return len(bytes), nil
}

func (m *ResponseWriter) WriteHeader(status int) {
	m.MockStatusCode = status
}

func (m *ResponseWriter) Flush() {
	// Noop.
}

type RequestBody struct {
	MockBuffer []byte
}

func (b *RequestBody) Read(p []byte) (int, error) {
	if len(b.MockBuffer) == 0 {
		return 0, io.EOF
	}

	n := copy(p, b.MockBuffer)
	b.MockBuffer = make([]byte, 0)

	return n, nil
}

func (b *RequestBody) Close() error {
	// Noop.
	return nil
}

func NewClient() *client.Client {
	srv := server.New()

	conf := &client.Config{
		ErrorLog:   srv.ErrorLog,
		InfoLog:    srv.InfoLog,
		PublicRoot: srv.PublicRoot,
		Efs:        srv.Efs,
	}

	writer := &ResponseWriter{
		MockHeader: map[string][]string{},
		MockBytes:  make([]byte, 0),
	}

	request := &http.Request{
		Header: map[string][]string{},
		Body: &RequestBody{
			MockBuffer: make([]byte, 1024),
		},
	}

	return &client.Client{
		Writer:  writer,
		Request: request,
		Config:  conf,
		EventId: 1,
		Status:  200,
	}
}
