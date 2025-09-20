package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/cli/user"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/text"
)

func TestDownload(t *testing.T) {
	url := "https://github.com/razshare/frizzante/archive/refs/heads/main.zip"

	var err error

	var cache string
	if cache, err = user.FrizzanteCache(); err != nil {
		t.Fatal(err)
	}

	hash := text.Sha1(url)

	cached := filepath.Join(cache, hash+".zip")

	if err = os.RemoveAll(cached); err != nil {
		t.Fatal(err)
	}
	if err = os.RemoveAll("main.zip"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("main.zip") }()
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
