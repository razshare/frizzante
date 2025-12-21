package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func AirConfig(options AirConfigOptions) (err error) {
	var data []byte
	var conf map[string]any

	if files.IsFile(".air.toml") {
		if data, err = os.ReadFile(".air.toml"); err != nil {
			return
		}
	} else {
		if data, err = options.Efs.ReadFile("internal/project/.air.toml"); err != nil {
			return
		}
	}

	if _, err = toml.Decode(string(data), &conf); err != nil {
		return
	}

	var tags []string
	if tags, err = tags_.Parse(options.Tags); err != nil {
		return
	}

	tagsString := strings.Join(tags, ",")
	build := conf["build"].(map[string]any)
	build["bin"] = filepath.Join(".gen", "tmp", "main.exe")

	if len(tagsString) > 0 {
		build["cmd"] = fmt.Sprintf("go build -tags %s -o %s .", tagsString, build["bin"])
	} else {
		build["cmd"] = fmt.Sprintf("go build -o %s .", build["bin"])
	}

	if err = os.Remove(".air.toml"); err != nil {
		return
	}

	var file *os.File
	if file, err = os.Create(filepath.Join(".air.toml")); err != nil {
		return
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
			return
		}
	}()

	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(conf); err != nil {
		return
	}

	messages.Success(".air.toml generated")

	return
}
