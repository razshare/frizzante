package menus

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
)

func LoadPlugins(menu *Menu, directory string) (err error) {
	type Configuration struct {
		Description string `toml:"description"`
		Hidden      bool   `toml:"hidden"`
	}
	if files.IsDirectory(directory) {
		var fileNames []string
		if fileNames, err = files.ReadDirectory(directory); err != nil {
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
					if !messages.Command(messages.CommandOptions{
						Environment: os.Environ(),
						Program:     *app.Go,
						Args:        []string{"run", fmt.Sprintf(".%s%s", string(filepath.Separator), packageName)},
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
