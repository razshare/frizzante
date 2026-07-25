package generations

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Migration(options MigrationOptions) (err error) {
	type Configuration struct {
		Sql []struct {
			Schema string `yaml:"schema"`
		} `yaml:"sql"`
	}
	baseDirectory := filepath.Join("migrate")
	spin := spinners.New("creating migration file")
	go spinners.Start(spin)
	if !files.IsDirectory(filepath.Join(baseDirectory, "migrations")) {
		if err = os.MkdirAll(filepath.Join(baseDirectory, "migrations"), os.ModePerm); err != nil {
			spinners.Stop(spin)
			return
		}
	}
	data := []byte("-- migration: down\n\n-- migration: up\n")
	now := time.Now()
	format := now.Format("2006_01_02T15_04_05Z07_00")
	migrationFileName := filepath.Join(baseDirectory, "migrations", fmt.Sprintf("%s.sql", format))
	if err = os.WriteFile(migrationFileName, data, os.ModePerm); err != nil {
		spinners.Stop(spin)
		return
	}
	spinners.Stop(spin)
	messages.Successf("migration generated at %s", migrationFileName)
	return
}
