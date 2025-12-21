package generate

import (
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
)

func Security(options SecurityOptions) (err error) {
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
