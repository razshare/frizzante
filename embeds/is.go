package embeds

import (
	"embed"
	"io/fs"
	"os"
)

// IsFile check if file exists and is a file.
func IsFile(efs embed.FS, name string) bool {
	var file fs.File
	var err error
	if file, err = efs.Open(name); err != nil {
		return false
	}

	var info os.FileInfo
	if info, err = file.Stat(); err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(efs embed.FS, n string) bool {
	file, err := efs.Open(n)
	if err != nil {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()
}
