package main

import (
	"fmt"
	"github.com/razshare/frizzante/globals"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestRenderServer(test *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	expected := "<h1>Welcome to Frizzante.</h1>"
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRenderServer", port))
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

func BenchmarkRenderServer(b *testing.B) {
	lock := <-server
	defer func() { server <- lock }()

	fmt.Println("=================================")
	fmt.Println("starting benchmark")
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	count := 1000
	var group sync.WaitGroup
	group.Add(count)

	for j := 0; j < count; j++ {
		go func() {
			defer group.Done()
			response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
			if getError != nil {
				b.Error(getError)
				return
			}

			_, readAllError := io.ReadAll(response.Body)
			if readAllError != nil {
				b.Error(readAllError)
				return
			}
		}()
	}

	group.Wait()

	runtime.ReadMemStats(&m2)
	fmt.Printf("used %d MB of memory\n", (m2.TotalAlloc-m1.TotalAlloc)/globals.MB)
	fmt.Printf("allocated memory %d times\n", m2.Mallocs-m1.Mallocs)
}

func TestRenderClient(test *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	expected := "<script type=\"application/javascript\">function target(){return document.getElementById("
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRenderClient", port))
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
