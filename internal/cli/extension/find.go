package extension

import "path/filepath"

func Find() string {
	if string(filepath.Separator) == "\\" {
		return ".exe"
	}

	return ""
}
