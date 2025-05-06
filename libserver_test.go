package frizzante

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestServerCreate(test *testing.T) {
	ServerCreate()
}

func TestServerWithHostName(test *testing.T) {
	server := ServerCreate()
	expected := "127.0.0.1"
	ServerWithHostName(server, expected)
	actual := server.hostName
	if actual != expected {
		test.Fatalf("server was expected to have host name '%s', received '%s' instead", expected, actual)
	}
}

func TestServerWithPort(test *testing.T) {
	server := ServerCreate()
	expected := 80
	ServerWithPort(server, expected)
	actual := server.port
	if actual != expected {
		test.Fatalf("server was expected to have port name %d, received %d instead", expected, actual)
	}
}

func TestServerWithReadTimeout(test *testing.T) {
	server := ServerCreate()
	expected := 10 * time.Second
	ServerWithReadTimeout(server, expected)
	actual := server.readTimeout
	if actual != expected {
		test.Fatalf("server was expected to have read timeout '%d', received '%d' instead", expected, actual)
	}

}

func TestServerWithWriteTimeout(test *testing.T) {
	server := ServerCreate()
	expected := 10 * time.Second
	ServerWithWriteTimeout(server, expected)
	actual := server.writeTimeout
	if actual != expected {
		test.Fatalf("server was expected to have write timeout '%d', received '%d' instead", expected, actual)
	}
}

func TestServerWithMaxHeaderBytes(test *testing.T) {
	server := ServerCreate()
	expected := 1 * MB
	ServerWithMaxHeaderBytes(server, expected)
	actual := server.maxHeaderBytes
	if actual != expected {
		test.Fatalf("server was expected to have max header bytes '%d', received '%d' instead", expected, actual)
	}
}

func TestServerWithCertificateAndKey(test *testing.T) {
	server := ServerCreate()
	expectedCertificate := "certificate.crt"
	expectedCertificateKey := "certificate.key"
	ServerWithCertificateAndKey(server, expectedCertificate, expectedCertificateKey)
	actualCertificate := server.certificate
	if actualCertificate != expectedCertificate {
		test.Fatalf("server was expected to have certificate '%s', received '%s' instead", expectedCertificate, actualCertificate)
	}
	actualCertificateKey := server.certificateKey
	if actualCertificateKey != expectedCertificateKey {
		test.Fatalf("server was expected to have certificate key '%s', received '%s' instead", expectedCertificateKey, actualCertificateKey)
	}
}

func TestServerWithTemporaryDirectory(test *testing.T) {
	server := ServerCreate()
	expected := ".temp"
	ServerWithTemporaryDirectory(server, expected)
	actual := server.temporaryDirectory
	if actual != expected {
		test.Fatalf("server was expected to have temporary directory '%s', received '%s' instead", expected, actual)
	}
}

func TestServerWithEmbeddedFileSystem(test *testing.T) {
	server := ServerCreate()
	expected := embeddedFileSystem
	ServerWithEmbeddedFileSystem(server, expected)
	actual := server.embeddedFileSystem
	if actual != expected {
		test.Fatalf("incorrect embedded file system detected")
	}
}

func TestServerTemporaryFileSave(test *testing.T) {
	server := ServerCreate()
	expected := "content"
	ServerWithTemporaryDirectory(server, ".temp")
	ServerTemporaryFileSave(server, "test", expected)
	fileName := filepath.Join(".temp", "test")
	if !Exists(fileName) {
		test.Fatalf("server was expected to create a temporary '%s', but it failed to do so", fileName)
	}
	bytes, readError := os.ReadFile(fileName)
	if readError != nil {
		test.Fatal(readError)
	}
	actual := string(bytes)
	if actual != expected {
		test.Fatalf("server temporary file was expected to contain '%s', received '%s' instead", expected, actual)
	}
}

func TestServerTemporaryFile(test *testing.T) {
	server := ServerCreate()
	expected := "content"
	ServerWithTemporaryDirectory(server, ".temp")
	ServerTemporaryFileSave(server, "test", expected)
	actual := ServerTemporaryFile(server, "test")
	if actual != expected {
		test.Fatalf("server temporary file was expected to contain '%s', received '%s' instead", expected, actual)
	}
}

func TestServerTemporaryFileExists(test *testing.T) {
	server := ServerCreate()
	ServerWithTemporaryDirectory(server, ".temp")
	ServerTemporaryFileSave(server, "test", "test")
	expected := true
	actual := ServerTemporaryFileExists(server, "test")
	if actual != expected {
		test.Fatalf("server was expected to have a temporary file by the name of 'test'")
	}
}

func TestServerTemporaryDirectoryClear(test *testing.T) {
	server := ServerCreate()
	ServerWithTemporaryDirectory(server, ".temp")
	ServerTemporaryFileSave(server, "test", "test")
	ServerTemporaryDirectoryClear(server)
	expected := false
	actual := ServerTemporaryFileExists(server, "test")
	if actual != expected {
		test.Fatalf("server was expected to not have a temporary file by the name of 'test'")
	}
}

func TestServerWithApi(test *testing.T) {
	server := ServerCreate()
	notifier := NotifierCreate()
	port := NextNumber(8080)
	ServerWithPort(server, port)
	ServerWithNotifier(server, notifier)
	expected := "hello"
	ServerWithApiBuilder(server, func(context ApiContext) {
		pattern, handler := context()
		pattern("GET /")
		handler(func(_ *Request, response *Response) {
			ResponseSendMessage(response, expected)
		})
	})
	go ServerStart(server)
	defer ServerStop(server)

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
	server := ServerCreate()
	notifier := NotifierCreate()
	port := NextNumber(8080)
	ServerWithPort(server, port)
	ServerWithNotifier(server, notifier)
	ServerWithApiBuilder(server, func(context ApiContext) {
		pattern, handler := context()
		pattern("GET /")
		handler(func(_ *Request, response *Response) {
			ResponseSendStatus(response, expected)
			ResponseSendMessage(response, "Ok")
		})
	})
	go ServerStart(server)
	defer ServerStop(server)

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
	server := ServerCreate()
	port := NextNumber(8080)
	notifier := NotifierCreate()
	ServerWithPort(server, port)
	ServerWithNotifier(server, notifier)
	expected := "application/json"
	ServerWithApiBuilder(server, func(context ApiContext) {
		pattern, handler := context()
		pattern("GET /")
		handler(func(_ *Request, response *Response) {
			ResponseSendHeader(response, "Content-Type", expected)
			ResponseSendMessage(response, "{}")
		})
	})
	go ServerStart(server)
	defer ServerStop(server)

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
