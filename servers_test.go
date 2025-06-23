package frizzante

import (
	"fmt"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/nums"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/servers"
	"io"
	ghttp "net/http"
	"testing"
	"time"
)

func TestNewServer(test *testing.T) {
	servers.New()
}

func TestServer_WithAddress(test *testing.T) {
	server := servers.New()
	expected := "127.0.0.1:8080"
	server.Address = "127.0.0.1:8080"
	actual := server.Address
	if actual != expected {
		test.Fatalf("server was expected to have host name '%s', received '%s' instead", expected, actual)
	}
}

func TestServer_WithReadTimeout(test *testing.T) {
	server := servers.New()
	expected := 10 * time.Second
	server.ReadTimeout = expected
	actual := server.ReadTimeout
	if actual != expected {
		test.Fatalf("server was expected to have read timeout '%d', received '%d' instead", expected, actual)
	}

}

func TestServer_WithWriteTimeout(test *testing.T) {
	server := servers.New()
	expected := 10 * time.Second
	server.WriteTimeout = expected
	actual := server.WriteTimeout
	if actual != expected {
		test.Fatalf("server was expected to have write timeout '%d', received '%d' instead", expected, actual)
	}
}

func TestServer_WithHeaderMaxMemory(test *testing.T) {
	server := servers.New()
	expected := 1 * globals.MB
	server.HeaderMaxMemory = expected
	actual := server.HeaderMaxMemory
	if actual != expected {
		test.Fatalf("server was expected to have max header bytes '%d', received '%d' instead", expected, actual)
	}
}

func TestServer_WithCertificate(test *testing.T) {
	server := servers.New()
	expectedCertificate := "cert.pem"
	expectedCertificateKey := "key.pem"
	server.Certificate = expectedCertificate
	server.Key = expectedCertificateKey
	actualCertificate := server.Certificate
	if actualCertificate != expectedCertificate {
		test.Fatalf("server was expected to have certificate '%s', received '%s' instead", expectedCertificate, actualCertificate)
	}
	actualCertificateKey := server.Key
	if actualCertificateKey != expectedCertificateKey {
		test.Fatalf("server was expected to have certificate key '%s', received '%s' instead", expectedCertificateKey, actualCertificateKey)
	}
}

func TestServer_AddRoute(test *testing.T) {
	expected := "hello"
	port := nums.NextNumber(8080)
	server := servers.New()
	server.Efs = emb
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(routes.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		con.SendMessage(expected)
	}})

	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	response, getError := ghttp.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
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
