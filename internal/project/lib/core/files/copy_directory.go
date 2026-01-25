package files

import (
	"path/filepath"
	"strings"
)

func CopyDirectory(from string, to string) (err error) {
	var entries []string
	if entries, err = ReadDirectory(from); err != nil {
		return
	}
	for _, entry := range entries {
		name := filepath.Join(to, strings.TrimPrefix(entry, from))
		if err = CopyFile(entry, name); err != nil {
			return
		}
	}
	return
}
