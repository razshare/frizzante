package generate

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
var TestDownloadEfs embed.FS

func TestDownload(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var file fs.File
		var err error
		if file, err = TestDownloadEfs.Open("test.zip"); err != nil {
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

	cached := filepath.Join(cache, hash)

	if err = os.RemoveAll(cached); err != nil {
		t.Fatal(err)
	}
	if err = os.RemoveAll("test_out.zip"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("test_out.zip") }()
	defer func() { _ = os.RemoveAll(cached) }()

	install, evict, err := Download(DownloadOptions{Auto: true, Url: testServer.URL})
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(cached) {
		t.Fatal("file should be cached")
	}

	installed, err := install(filepath.Join(".gen", "download"))
	if err != nil {
		t.Fatal(err)
	}

	if !installed {
		t.Fatal("resource should be installed in .gen/download")
	}

	if !files.IsDirectory(".gen/download") {
		t.Fatal(".gen/download should exist")
	}

	if err = evict(); err != nil {
		t.Fatal(err)
	}

	if files.IsFile(cached) {
		t.Fatal("file cache should be evicted")
	}
}
