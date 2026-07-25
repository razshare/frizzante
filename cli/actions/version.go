package actions

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/v2/tui/configs"
)

func Version(options VersionOptions) (err error) {
	var data []byte
	if data, err = options.Efs.ReadFile("version"); err != nil {
		return
	}
	version := string(data)
	lines := strings.Split(version, "\n")
	if len(lines) == 0 {
		return
	}
	fmt.Println(configs.Styles.Menu.PaddingRight(1).Render("│") + lines[0])
	return
}
