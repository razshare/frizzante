package actions

import (
	"fmt"
	"os"
	"strings"

	"github.com/razshare/frizzante/cli/tags"
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
				args := []string{"run"}
				if len(options.Tags) > 0 {
					var parsedTags []string
					if parsedTags, err = tags.Parse(options.Tags); err != nil {
						return
					}
					tagsLocal := make([]string, 0)
					for _, parsedTag := range parsedTags {
						tagsLocal = append(tagsLocal, parsedTag)
					}
					args = append(args, fmt.Sprintf("-tags=%s", strings.Join(tagsLocal, ",")))
				}
				args = append(args, "./post")
				messages.Command(messages.CommandOptions{
					Environment: os.Environ(),
					Program:     options.Go,
					Args:        args,
				})
				break
			}
		}
	}
	return
}
