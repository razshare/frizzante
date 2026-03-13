package menus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
)

func LoadPlugins(menu *Menu, directoryName string) (err error) {
	type Configuration struct {
		Description string `toml:"description"`
		Hidden      bool   `toml:"hidden"`
	}
	if files.IsDirectory(directoryName) {
		var fileNames []string
		if fileNames, err = files.ReadDirectory(directoryName); err != nil {
			return
		}
		for _, fileName := range fileNames {
			if filepath.Base(fileName) != "main.toml" {
				continue
			}
			var data []byte
			if data, err = os.ReadFile(fileName); err != nil {
				return
			}
			var configuration Configuration
			if err = toml.Unmarshal(data, &configuration); err != nil {
				return
			}
			packageName := filepath.Dir(fileName)
			id := filepath.Base(packageName)
			menu.Items = append(menu.Items, Item{
				Hidden: configuration.Hidden,
				Choice: search.Choice{Id: id, Description: configuration.Description},
				Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == id },
				Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
					parameters := make([]string, 0)
					if *app.Strict {
						parameters = append(parameters, "--strict")
					}
					if *app.Incremental {
						parameters = append(parameters, "--incremental")
					}
					if *app.Bun != "" {
						parameters = append(parameters, fmt.Sprintf("--bun=\"%s\"", *app.Bun))
					}
					if *app.Go != "" {
						parameters = append(parameters, fmt.Sprintf("--go=\"%s\"", *app.Go))
					}
					if *app.Tags != "" {
						parameters = append(parameters, fmt.Sprintf("--tags=\"%s\"", *app.Tags))
					}
					if *app.Air != "" {
						parameters = append(parameters, fmt.Sprintf("--air=\"%s\"", *app.Air))
					}
					if *app.Database != "" {
						parameters = append(parameters, fmt.Sprintf("--database=\"%s\"", *app.Database))
					}
					if *app.DatabaseType != "" {
						parameters = append(parameters, fmt.Sprintf("--database-type=\"%s\"", *app.DatabaseType))
					}
					if *app.Sqlc != "" {
						parameters = append(parameters, fmt.Sprintf("--sqlc=\"%s\"", *app.Sqlc))
					}
					if *app.SqlcYaml != "" {
						parameters = append(parameters, fmt.Sprintf("--sqlc-yaml=\"%s\"", *app.SqlcYaml))
					}
					parameters = append(parameters, value)
					parameters = append(parameters, query...)
					packageDirectoryName := fmt.Sprintf(".%s%s", string(filepath.Separator), packageName)
					goRun := []string{"run", packageDirectoryName}
					if len(parameters) > 0 {
						goRun = append(goRun, strings.Join(parameters, " "))
					}
					if !messages.Command(messages.CommandOptions{
						Environment: os.Environ(),
						Program:     *app.Go,
						Args:        goRun,
					}) {
						err = fmt.Errorf("plugin %s failed to execute", packageName)
						return
					}
					messages.Successf("plugin %s executed successfully", packageName)
					return
				},
			})
		}
	}
	return
}
