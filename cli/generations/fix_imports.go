package generations

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
)

func FixImports(options FixImportsOptions) (err error) {
	if !files.IsDirectory(options.Directory) {
		err = fmt.Errorf("directory %s not found", options.Directory)
		return
	}
	befores := [][]byte{
		[]byte("github.com/razshare/frizzante/v2/internal/project"),
	}
	after := []byte("main")
	var entries []string
	if entries, err = files.ReadDirectory(options.Directory); err != nil {
		return
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry, ".go") &&
			!strings.HasSuffix(entry, ".svelte") &&
			!strings.HasSuffix(entry, ".js") &&
			!strings.HasSuffix(entry, ".ts") {
			continue
		}
		var data []byte
		if data, err = os.ReadFile(entry); err != nil {
			return
		}
		for _, before := range befores {
			data = bytes.ReplaceAll(data, before, after)
		}
		if err = os.WriteFile(entry, data, os.ModePerm); err != nil {
			return
		}
	}
	return
}
