package generations

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Types(options TypesOptions) (err error) {
	if !files.IsDirectory(filepath.Join("lib", "types")) {
		err = errors.New("types directory not found")
		return
	}
	program := "." + string(filepath.Separator) + filepath.Join("lib", "types")
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"run", program},
	}) {
		err = errors.New("could not generate type definitions")
		return
	}
	messages.Success("type definitions generated")
	return
}
