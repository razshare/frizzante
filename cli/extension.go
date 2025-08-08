package cli

import "path/filepath"

func Extension() string {
	if string(filepath.Separator) == "\\" {
		return ".exe"
	}

	return ""
}
