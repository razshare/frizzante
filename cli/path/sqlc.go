package path

import (
	"github.com/razshare/frizzante/cli/extension"
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/tui/messages"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Sqlc(basepath string) string {
	var sqlc string

	if *flags.Sqlc != "" {
		sqlc = *flags.Sqlc
	} else {
		sqlc = filepath.Join(".gen", "sqlc", "sqlc")
	}

	if strings.HasPrefix(sqlc, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		sqlc = strings.Replace(sqlc, "~", dirname, 1)
		return sqlc + extension.Find()
	}

	if !strings.Contains(sqlc, string(filepath.Separator)) {
		return sqlc + extension.Find()
	}

	var pathError error
	sqlc, pathError = filepath.Rel(basepath, sqlc)
	if pathError != nil {
		messages.Fatal(pathError)
	}

	return sqlc + extension.Find()
}
