package frizzante

import (
	"fmt"
	"github.com/razshare/frizzante/libcon"
	"github.com/razshare/frizzante/libnum"
	"github.com/razshare/frizzante/libsrv"
	"github.com/razshare/frizzante/libview"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRenderServer(test *testing.T) {
	port := libnum.NextNumber(8080)
	server := libsrv.NewServer()
	server.Efs = efs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(libsrv.Route{Pattern: "GET /welcome", Handler: func(c *libcon.Connection) {
		c.SendView(libview.View{
			Name:       "Welcome",
			RenderMode: libview.RenderModeServer,
			Data:       map[string]any{"name": "world"},
		})
	}})

	go server.Start()
	defer server.Stop()
	time.Sleep(1 * time.Second)

	expected := "<h1>Hello world.</h1>"
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/welcome", port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}

func TestRenderClient(test *testing.T) {
	port := libnum.NextNumber(8080)
	server := libsrv.NewServer()
	server.Efs = efs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(libsrv.Route{Pattern: "GET /welcome", Handler: func(c *libcon.Connection) {
		c.SendView(libview.View{
			Name:       "Welcome",
			RenderMode: libview.RenderModeClient,
			Data:       map[string]any{"name": "world"},
		})
	}})

	go server.Start()
	defer server.Stop()
	time.Sleep(1 * time.Second)

	expected := "<script type=\"application/javascript\">function target(){return document.getElementById("
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/welcome", port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}
