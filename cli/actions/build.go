package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Build(options BuildOptions) (err error) {
	if err = Package(PackageOptions{
		Go:  options.Go,
		Bun: options.Bun,
	}); err != nil {
		return
	}
	if err = PreBuild(PreBuildOptions(options)); err != nil {
		return
	}
	extension := extensions.Find()
	spin := spinners.New("building binary")
	args := []string{"build", fmt.Sprintf("-o=%s", filepath.Join(".gen", "bin", "app"+extension))}
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
	args = append(args, ".")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        args,
	}) {
		messages.Error("could not build go source code")
		return
	}
	if err = PostBuild(PostBuildOptions(options)); err != nil {
		return
	}
	spinners.Stop(spin)
	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"+extension))
	return
}
