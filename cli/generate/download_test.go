package generate

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/razshare/frizzante/cli/paths"
	"github.com/razshare/frizzante/internal/additions/lib/security"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
)

//go:embed test.zip
var TestDownloadEfs embed.FS

func TestDownload(t *testing.T) {
	server := servers.New()
	server.Addr = "127.0.0.1:7676"
	server.Routes = []routes.Route{
		{Pattern: "GET /", Handler: func(client *clients.Client) {
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

			http.ServeContent(client.Writer, &client.Request, "test_out.zip", info.ModTime(), bytes.NewReader(buf))
		}},
	}
	go servers.Start(server)
	defer func() { server.Channels.End <- struct{}{} }()

	time.Sleep(time.Second)

	url := "http://127.0.0.1:7676/test_out.zip"

	var err error

	var cache string
	if cache, err = paths.Cache(); err != nil {
		t.Fatal(err)
	}

	hash := security.Sha1(url)

	cached := filepath.Join(cache, hash+".zip")

	if err = os.RemoveAll(cached); err != nil {
		t.Fatal(err)
	}
	if err = os.RemoveAll("test_out.zip"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("test_out.zip") }()
	defer func() { _ = os.RemoveAll(cached) }()

	install, evict, err := Download(DownloadOptions{Auto: true, Url: url})
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
