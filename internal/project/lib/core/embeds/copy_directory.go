package embeds

import (
	"embed"
	"path/filepath"
	"strings"
)

func CopyDirectory(efs embed.FS, from string, to string) (err error) {
	var entries []string
	if entries, err = ReadDirectory(efs, from); err != nil {
		return
	}
	for _, ent := range entries {
		n := filepath.Join(to, strings.TrimPrefix(ent, from))
		err = CopyFile(efs, ent, n)
		if err != nil {
			return
		}
	}
	return nil
}
