package actions

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/razshare/frizzante/cli/extensions"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/tui/messages"
)

func Snapshot(options SnapshotOptions) (err error) {
	var tags []string
	if tags, err = tags_.Parse(options.Tags); err != nil {
		return
	}

	if !slices.Contains(tags, "snapshot_servers") {
		tags = append(tags, "snapshot_servers")
	}

	if err = Build(BuildOptions{
		Go:   options.Go,
		Tags: strings.Join(tags, ","),
		Bun:  options.Bun,
	}); err != nil {
		return
	}

	extension := extensions.Find()

	messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     filepath.Join(".gen", "bin", "app"+extension),
	})

	return
}
