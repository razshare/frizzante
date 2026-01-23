package mocks

import (
	"io"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/ssr"
)

type ResponseWriter struct {
	MockHeader     http.Header
	MockStatusCode int
	MockBytes      []byte
}

func (model *ResponseWriter) Header() http.Header {
	return model.MockHeader
}

func (model *ResponseWriter) Write(bytes []byte) (int, error) {
	model.MockBytes = append(model.MockBytes, bytes...)
	return len(bytes), nil
}

func (model *ResponseWriter) WriteHeader(status int) {
	model.MockStatusCode = status
}

func (model *ResponseWriter) Flush() {
	// Noop.
}

type RequestBody struct {
	MockBuffer []byte
}

func (body *RequestBody) Read(p []byte) (int, error) {
	if len(body.MockBuffer) == 0 {
		return 0, io.EOF
	}

	count := copy(p, body.MockBuffer)
	body.MockBuffer = make([]byte, 0)

	return count, nil
}

func (body *RequestBody) Close() error {
	// Noop.
	return nil
}

func NewClient() *clients.Client {
	srv := ssr.NewServer()

	conf := clients.Options{
		ErrorLog: srv.ErrorLog,
		InfoLog:  srv.InfoLog,
		Efs:      srv.Efs,
	}

	writer := &ResponseWriter{
		MockHeader: map[string][]string{},
		MockBytes:  make([]byte, 0),
	}

	request := http.Request{
		Header: map[string][]string{},
		Body: &RequestBody{
			MockBuffer: make([]byte, 1024),
		},
	}

	return &clients.Client{
		Writer:  writer,
		Request: request,
		Options: conf,
		EventId: 1,
		Status:  200,
	}
}
