package actions

import (
	"os"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func PostBuild(options PostBuildOptions) (err error) {
	if files.IsDirectory("post") {
		var fileNames []string
		if fileNames, err = files.ReadDirectory("post"); err != nil {
			return
		}
		for _, fileName := range fileNames {
			if strings.HasSuffix(fileName, ".go") {
				messages.Command(messages.CommandOptions{
					Environment: os.Environ(),
					Program:     options.Go,
					Args:        []string{"run", "./post"},
				})
				break
			}
		}
	}
	return
}
