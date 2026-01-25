package paths

import (
	"os"
	"path/filepath"
)

func Home() (home string, err error) {
	if home = os.Getenv("FRIZZANTE_HOME"); home == "" {
		var user string
		if user, err = os.UserHomeDir(); err != nil {
			return "", err
		}
		home = filepath.Join(user, ".frizzante")
	}
	return
}
