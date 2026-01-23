package generations

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Types(options TypesOptions) (err error) {
	if !files.IsDirectory("types") {
		err = errors.New("types directory not found")
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"run", "./types"},
	}) {
		err = errors.New("could not generate type definitions")
		return
	}
	messages.Success("type definitions generated")
	return
}
