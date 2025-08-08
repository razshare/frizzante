package embeds

import "embed"

// IsFile check if file exists and is a file.
func IsFile(efs embed.FS, n string) bool {
	file, err := efs.Open(n)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory checks if file exists and is a directory.
func IsDirectory(efs embed.FS, n string) bool {
	file, err := efs.Open(n)
	if err != nil {
		return false
	}
	stat, statError := file.Stat()
	if statError != nil {
		return false
	}
	return stat.IsDir()
}
