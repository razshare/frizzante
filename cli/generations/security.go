package generations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
)

func Security(options SecurityOptions) (err error) {
	directoryName := filepath.Join("lib", "security")
	if files.IsDirectory(directoryName) {
		if options.Strict {
			err = fmt.Errorf("%s already exists", directoryName)
			return
		}
		var yesRemove bool
		if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
			return err
		}
		if yesRemove {
			if err = os.RemoveAll(directoryName); err != nil {
				return
			}
		}
	}

	if err = Copy(CopyOptions{
		From: "internal/additions/lib/security",
		To:   filepath.Join("lib", "security"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "security")}); err != nil {
		return
	}

	messages.Command(messages.CommandOptions{Environment: os.Environ(), Program: "go", Args: []string{"get", "golang.org/x/crypto/bcrypt"}})
	messages.Command(messages.CommandOptions{Environment: os.Environ(), Program: "go", Args: []string{"get", "golang.org/x/crypto/sha3"}})
	messages.Command(messages.CommandOptions{Environment: os.Environ(), Program: "go", Args: []string{"get", "golang.org/x/text/unicode/norm"}})

	return
}
