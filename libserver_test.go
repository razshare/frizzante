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

func TestServerWithAddress(test *testing.T) {
	server := NewServer()
	expected := "127.0.0.1:8080"
	server.WithAddress("127.0.0.1:8080")
	actual := server.address
	if actual != expected {
		test.Fatalf("server was expected to have host name '%s', received '%s' instead", expected, actual)
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
	server.WithHeaderMaxMemory(expected)
	actual := server.headerMaxMemory
	if actual != expected {
		test.Fatalf("server was expected to have max header bytes '%d', received '%d' instead", expected, actual)
	}
}

func TestServerWithCertificateAndKey(test *testing.T) {
	server := NewServer()
	expectedCertificate := "cert.pem"
	expectedCertificateKey := "key.pem"
	server.WithCertificate(expectedCertificate, expectedCertificateKey)
	actualCertificate := server.certificate
	if actualCertificate != expectedCertificate {
		test.Fatalf("server was expected to have certificate '%s', received '%s' instead", expectedCertificate, actualCertificate)
	}
	actualCertificateKey := server.key
	if actualCertificateKey != expectedCertificateKey {
		test.Fatalf("server was expected to have certificate key '%s', received '%s' instead", expectedCertificateKey, actualCertificateKey)
	}
}

func TestServerWithApi(test *testing.T) {
	server := NewServer()
	notifier := NewNotifier()
	port := NextNumber(8080)
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.WithNotifier(notifier)
	expected := "hello"
	server.WithRoute("GET /", []Guard{}, func(request *Request, response *Response) {
		response.SendMessage(expected)
	})
	go server.Start(efs)
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
	port := NextNumber(8080)
	notifier := NewNotifier()
	server := NewServer()
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.WithNotifier(notifier)
	server.WithRoute("GET /", []Guard{}, func(request *Request, response *Response) {
		response.SendStatus(expected)
		response.SendMessage("ok")
	})
	go server.Start(efs)
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
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.WithNotifier(notifier)
	expected := "application/json"
	server.WithRoute("GET /", []Guard{}, func(req *Request, res *Response) {
		res.SendHeader("Content-Type", expected)
		res.SendMessage("{}")
	})
	go server.Start(efs)
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
