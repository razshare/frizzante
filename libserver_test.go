package frizzante

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(test *testing.T) {
	NewServer()
}

func TestServerWithHostName(test *testing.T) {
	server := NewServer()
	expected := "127.0.0.1"
	server.WithHostName(expected)
	actual := server.hostName
	if actual != expected {
		test.Fatalf("server was expected to have host name '%s', received '%s' instead", expected, actual)
	}
}

func TestServerWithPort(test *testing.T) {
	server := NewServer()
	expected := 80
	server.WithPort(expected)
	actual := server.port
	if actual != expected {
		test.Fatalf("server was expected to have port name %d, received %d instead", expected, actual)
	}
}

func TestServerWithReadTimeout(test *testing.T) {
	server := NewServer()
	expected := 10 * time.Second
	server.WithReadTimeout(expected)
	actual := server.readTimeout
	if actual != expected {
		test.Fatalf("server was expected to have read timeout '%d', received '%d' instead", expected, actual)
	}

}

func TestServerWithWriteTimeout(test *testing.T) {
	server := NewServer()
	expected := 10 * time.Second
	server.WithWriteTimeout(expected)
	actual := server.writeTimeout
	if actual != expected {
		test.Fatalf("server was expected to have write timeout '%d', received '%d' instead", expected, actual)
	}
}

func TestServerWithMaxHeaderBytes(test *testing.T) {
	server := NewServer()
	expected := 1 * MB
	server.WithMaxHeaderBytes(expected)
	actual := server.maxHeaderBytes
	if actual != expected {
		test.Fatalf("server was expected to have max header bytes '%d', received '%d' instead", expected, actual)
	}
}

func TestServerWithCertificateAndKey(test *testing.T) {
	server := NewServer()
	expectedCertificate := "certificate.crt"
	expectedCertificateKey := "certificate.key"
	server.WithCertificateAndKey(expectedCertificate, expectedCertificateKey)
	actualCertificate := server.certificate
	if actualCertificate != expectedCertificate {
		test.Fatalf("server was expected to have certificate '%s', received '%s' instead", expectedCertificate, actualCertificate)
	}
	actualCertificateKey := server.certificateKey
	if actualCertificateKey != expectedCertificateKey {
		test.Fatalf("server was expected to have certificate key '%s', received '%s' instead", expectedCertificateKey, actualCertificateKey)
	}
}

func TestServerWithEmbeddedFileSystem(test *testing.T) {
	server := NewServer()
	expected := &embeddedFileSystem
	server.WithEmbeddedFileSystem(expected)
	actual := server.embeddedFileSystem
	if actual != expected {
		test.Fatalf("incorrect embedded file system detected")
	}
}

func TestServerWithApi(test *testing.T) {
	server := NewServer()
	notifier := NewNotifier()
	port := NextNumber(8080)
	server.WithPort(port)
	server.WithNotifier(notifier)
	expected := "hello"
	server.OnRequest("GET /", func(request *Request, response *Response) {
		response.SendMessage(expected)
	})
	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	actual, getError := HttpGet(fmt.Sprintf("http://127.0.0.1:%d/", port), nil)
	if getError != nil {
		test.Fatal(getError)
	}

	if actual != expected {
		test.Fatalf("server was expected to respond with '%s', received '%s' instead", expected, actual)
	}
}

func TestSendStatus(test *testing.T) {
	expected := 201
	server := NewServer()
	notifier := NewNotifier()
	port := NextNumber(8080)
	server.WithPort(port)
	server.WithNotifier(notifier)
	server.OnRequest("GET /", func(request *Request, response *Response) {
		response.SendStatus(expected)
		response.SendMessage("ok")
	})
	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer response.Body.Close()

	actual := response.StatusCode

	if actual != expected {
		test.Fatalf("server was expected to respond with status code '%d', received '%d' intead", expected, actual)
	}
}

func TestSendHeader(test *testing.T) {
	server := NewServer()
	port := NextNumber(8080)
	notifier := NewNotifier()
	server.WithPort(port)
	server.WithNotifier(notifier)
	expected := "application/json"
	server.OnRequest("GET /", func(req *Request, res *Response) {
		res.SendHeader("Content-Type", expected)
		res.SendMessage("{}")
	})
	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer response.Body.Close()

	actual := response.Header.Get("Content-Type")

	if actual != expected {
		test.Fatalf("server was expected to respond with header content type '%s', received '%s' intead", expected, actual)
	}
}
