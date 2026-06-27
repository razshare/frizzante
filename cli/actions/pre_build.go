package actions

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func PreBuild(options PreBuildOptions) (err error) {
	if !files.IsDirectory("pre") {
		return
	}
	var fileNames []string
	if fileNames, err = files.FindWithSuffix("pre", ".go"); err != nil {
		return
	}
	if len(fileNames) == 0 {
		return
	}
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
	args = append(args, fmt.Sprintf(".%spre", string(filepath.Separator)))
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        args,
	}) {
		err = errors.New("something went wrong while running prebuild")
		return
	}
	return
}
