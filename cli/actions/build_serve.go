package actions

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func BuildServe(options BuildServeOptions) (err error) {
	if err = PreBuild(PreBuildOptions{
		Go:   options.Go,
		Tags: options.Tags,
	}); err != nil {
		return
	}
	if err = Package(PackageOptions{
		Go:  options.Go,
		Bun: options.Bun,
	}); err != nil {
		return
	}
	output := filepath.Join(strings.ReplaceAll(options.Output, "/", string(filepath.Separator)), "serve")
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
	buildAppArgs := []string{"build", fmt.Sprintf("-o=%s", output)}
	buildAppArgs = append(buildAppArgs, buildTags...)
	buildAppArgs = append(buildAppArgs, ".")
	spin := spinners.New("building serve binary")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        buildAppArgs,
	}) {
		spinners.Stop(spin)
		err = errors.New("something went wrong while building serve")
		return
	}
	spinners.Stop(spin)
	messages.Success("project serve built into ", output)
	if err = BuildMigrate(BuildMigrateOptions{
		Go:     options.Go,
		Tags:   options.Tags,
		Output: options.Output,
	}); err != nil {
		return
	}
	if err = PostBuild(PostBuildOptions{
		Go:   options.Go,
		Tags: options.Tags,
	}); err != nil {
		return
	}
	return
}
