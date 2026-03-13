package embeds

import (
	"embed"
	"fmt"
	"io/fs"
	"slices"
)

func ReadDirectory(efs embed.FS, from string) (entries []string, err error) {
	entries = make([]string, 0)
	var names []fs.DirEntry
	if names, err = efs.ReadDir(from); err != nil {
		return
	}
	for _, name := range names {
		if name.IsDir() {
			var subentries []string
			subentries, err = ReadDirectory(efs, fmt.Sprintf("%s/%s", from, name.Name()))
			if err != nil {
				return
			}
			entries = slices.Concat(entries, subentries)
			continue
		}
		entries = append(entries, fmt.Sprintf("%s/%s", from, name.Name()))
	}
	return
}
