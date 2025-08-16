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

func TestRenderServer(t *testing.T) {
	<-serve
	defer func() { serve <- 0 }()

	ex := "<h1>Welcome to Frizzante.</h1>"
	r, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRenderServer", port))
	if err != nil {
		t.Fatal(err)
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	ac := string(d)

	ok := strings.Contains(ac, ex)

	if !ok {
		t.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", ex, ac)
	}
}

func BenchmarkRenderServer(b *testing.B) {
	<-serve
	defer func() { serve <- 0 }()

	count := 1000

	fmt.Println("=================================")
	fmt.Printf("starting benchmark, attempting to send %d requests, all in parallel\n", count)
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	var grp sync.WaitGroup
	grp.Add(count)

	for j := 0; j < count; j++ {
		go func() {
			defer grp.Done()
			r, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
			if err != nil {
				b.Error(err)
				return
			}

			_, err = io.ReadAll(r.Body)
			if err != nil {
				b.Error(err)
				return
			}
		}()
	}

	grp.Wait()

	runtime.ReadMemStats(&m2)
	fmt.Printf("used %d MB of memory\n", (m2.TotalAlloc-m1.TotalAlloc)/globals.MB)
	fmt.Printf("allocated memory %d times\n", m2.Mallocs-m1.Mallocs)
}

func TestRenderClient(t *testing.T) {
	<-serve
	defer func() { serve <- 0 }()

	ex := "return document.getElementById(\"svelte-app\")"
	r, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRenderClient", port))
	if err != nil {
		t.Fatal(err)
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	ac := string(d)

	ok := strings.Contains(ac, ex)

	if !ok {
		t.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", ex, ac)
	}
}
