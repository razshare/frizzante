package frz

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestNewServer(test *testing.T) {
	NewServer()
}

func TestServer_WithAddress(test *testing.T) {
	server := NewServer()
	expected := "127.0.0.1:8080"
	server.WithAddress("127.0.0.1:8080")
	actual := server.address
	if actual != expected {
		test.Fatalf("server was expected to have host name '%s', received '%s' instead", expected, actual)
	}
}

func TestServer_WithReadTimeout(test *testing.T) {
	server := NewServer()
	expected := 10 * time.Second
	server.WithReadTimeout(expected)
	actual := server.readTimeout
	if actual != expected {
		test.Fatalf("server was expected to have read timeout '%d', received '%d' instead", expected, actual)
	}

}

func TestServer_WithWriteTimeout(test *testing.T) {
	server := NewServer()
	expected := 10 * time.Second
	server.WithWriteTimeout(expected)
	actual := server.writeTimeout
	if actual != expected {
		test.Fatalf("server was expected to have write timeout '%d', received '%d' instead", expected, actual)
	}
}

func TestServer_WithHeaderMaxMemory(test *testing.T) {
	server := NewServer()
	expected := 1 * MB
	server.WithHeaderMaxMemory(expected)
	actual := server.headerMaxMemory
	if actual != expected {
		test.Fatalf("server was expected to have max header bytes '%d', received '%d' instead", expected, actual)
	}
}

func TestServer_WithCertificate(test *testing.T) {
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

func TestServer_AddRoute(test *testing.T) {
	expected := "hello"
	port := NextNumber(8080)
	server := NewServer().
		WithEfs(dist).
		WithAddress(fmt.Sprintf("127.0.0.1:%d", port)).
		AddRoute(Route{Pattern: "GET /", Handler: func(c *Connection) {
			c.SendMessage(expected)
		}})

	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	if actual != expected {
		test.Fatalf("server was expected to respond with '%s', received '%s' instead", expected, actual)
	}
}
