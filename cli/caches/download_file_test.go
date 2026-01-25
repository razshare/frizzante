package caches

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/cli/paths"
	"github.com/razshare/frizzante/internal/additions/lib/security"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed test.zip
var TestDownloadFileEfs embed.FS

func TestDownloadFile(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var file fs.File
		var err error
		if file, err = TestDownloadFileEfs.Open("test.zip"); err != nil {
			t.Fatal(err)
		}
		var info os.FileInfo
		if info, err = file.Stat(); err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			t.Fatal(err)
		}
		http.ServeContent(writer, request, "test_out.zip", info.ModTime(), bytes.NewReader(buf))
	}))
	defer func() { testServer.Close() }()
	var err error
	var cache string
	if cache, err = paths.Cache(); err != nil {
		t.Fatal(err)
	}
	hash := security.Sha1(testServer.URL)
	cachedFileName := filepath.Join(cache, hash)
	if err = os.RemoveAll(cachedFileName); err != nil {
		t.Fatal(err)
	}
	if err = os.RemoveAll("test_out.zip"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("test_out.zip") }()
	defer func() { _ = os.RemoveAll(cachedFileName) }()
	if _, err = DownloadFile(DownloadFileOptions{
		Url: testServer.URL,
	}); err != nil {
		t.Fatal(err)
	}
	if !files.IsFile(cachedFileName) {
		t.Fatal("file should be cached")
	}
}
