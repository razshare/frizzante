package generations

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func AirConfig(options AirConfigOptions) (err error) {
	var data []byte
	var conf map[string]any
	airTomlExists := files.IsFile(".air.toml")
	if airTomlExists {
		if data, err = os.ReadFile(".air.toml"); err != nil {
			return
		}
	} else {
		if data, err = options.Efs.ReadFile("internal/project/.air.toml"); err != nil {
			return
		}
	}
	content := string(data)
	if _, err = toml.Decode(content, &conf); err != nil {
		return
	}
	if airTomlExists {
		if err = os.Remove(".air.toml"); err != nil {
			return
		}
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
