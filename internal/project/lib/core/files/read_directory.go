package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

func ReadDirectory(from string) (names []string, err error) {
	var entries []fs.DirEntry
	if entries, err = os.ReadDir(from); err != nil {
		return
	}
	for _, entry := range entries {
		name := from + string(filepath.Separator) + entry.Name()
		if entry.IsDir() {
			var subnames []string
			if subnames, err = ReadDirectory(name); err != nil {
				return
			}
			names = slices.Concat(names, subnames)
			continue
		}
		names = append(names, name)
	}
	return
}
