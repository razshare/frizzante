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
	"github.com/razshare/frizzante/tui/spinners"
)

func BuildMigrate(options BuildMigrateOptions) (err error) {
	if !files.IsDirectory("migrate") {
		return
	}
	var fileNames []string
	if fileNames, err = files.FindWithSuffix("migrate", ".go"); err != nil {
		return
	}
	if len(fileNames) == 0 {
		return
	}
	output := filepath.Join(strings.ReplaceAll(options.Output, "/", string(filepath.Separator)), "migrate")
	buildTags := make([]string, 0)
	if len(options.Tags) > 0 {
		var parsedTags []string
		if parsedTags, err = tags.Parse(options.Tags); err != nil {
			return
		}
		tagsLocal := make([]string, 0)
		for _, parsedTag := range parsedTags {
			tagsLocal = append(tagsLocal, parsedTag)
		}
		buildTags = append(buildTags, fmt.Sprintf("-tags=%s", strings.Join(tagsLocal, ",")))
	}
	spin := spinners.New("building migrate binary")
	go spinners.Start(spin)
	buildMigrateArgs := []string{"build", fmt.Sprintf("-o=%s", output)}
	buildMigrateArgs = append(buildMigrateArgs, buildTags...)
	buildMigrateArgs = append(buildMigrateArgs, fmt.Sprintf(".%smigrate", string(filepath.Separator)))
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        buildMigrateArgs,
	}) {
		spinners.Stop(spin)
		err = errors.New("something went wrong while building migrate")
		return
	}
	spinners.Stop(spin)
	messages.Success("project migrate built into ", output)
	return
}
