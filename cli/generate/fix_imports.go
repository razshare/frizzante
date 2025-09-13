package generate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func FixImports(options FixImportsOptions) (err error) {
	before := []byte("github.com/razshare/frizzante/internal/project")
	after := []byte("main")

	var entries []string
	if entries, err = files.ReadDirectory(options.Directory); err != nil {
		return
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry, ".go") {
			continue
		}

		if entry != "main.go" && strings.HasPrefix(entry, "lib"+string(filepath.Separator)) {
			continue
		}

		var data []byte
		if data, err = os.ReadFile(entry); err != nil {
			return
		}

		data = bytes.ReplaceAll(data, before, after)

		if err = os.WriteFile(entry, data, os.ModePerm); err != nil {
			return
		}
	}

	return
}
