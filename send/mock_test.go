package send

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/server"
	"net/http"
)

type MockWriter struct {
	MockHeader     http.Header
	MockStatusCode int
	MockBytes      []byte
}

func (m *MockWriter) Header() http.Header {
	return m.MockHeader
}

func (m *MockWriter) Write(bytes []byte) (int, error) {
	m.MockBytes = append(m.MockBytes, bytes...)
	return len(bytes), nil
}

func (m *MockWriter) WriteHeader(status int) {
	m.MockStatusCode = status
}

func (m *MockWriter) Flush() {
	// Noop.
}

func MockClient() *client.Client {
	srv := server.New()

	conf := &client.Config{
		Render:     srv.Render,
		ErrorLog:   srv.ErrorLog,
		InfoLog:    srv.InfoLog,
		PublicRoot: srv.PublicRoot,
		Efs:        srv.Efs,
	}

	writer := &MockWriter{
		MockHeader: map[string][]string{},
		MockBytes:  make([]byte, 0),
	}

	request := &http.Request{
		Header: map[string][]string{},
	}

	return &client.Client{
		Writer:  writer,
		Request: request,
		Config:  conf,
		EventId: 1,
		Status:  200,
	}
}
