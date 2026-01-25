package embeds

import "embed"

// IsDirectory checks if file exists and is a directory.
func IsDirectory(efs embed.FS, name string) bool {
	file, err := efs.Open(name)
	if err != nil {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()
}
