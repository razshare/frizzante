package on

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Build(app string, plat string, gobin string, bunbin string) error {
	err := Package(app, bunbin)
	if err != nil {
		return err
	}

	build := exec.Command(gobin, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".")
	build.Env = os.Environ()

	if strings.ToLower(plat) == "linux/amd64" {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=amd64")
	} else if strings.ToLower(plat) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=arm64")
	} else if strings.ToLower(plat) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=amd64")
	} else if strings.ToLower(plat) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=arm64")
	}

	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Stdin = os.Stdin
	err = build.Run()
	if err != nil {
		return err
	}

	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))

	return nil
}
