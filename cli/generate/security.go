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
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "security")}); err != nil {
		return
	}

	messages.Command("", os.Environ(), "go", "get", "golang.org/x/crypto/bcrypt")
	messages.Command("", os.Environ(), "go", "get", "golang.org/x/crypto/sha3")
	messages.Command("", os.Environ(), "go", "get", "golang.org/x/text/unicode/norm")

	return
}
