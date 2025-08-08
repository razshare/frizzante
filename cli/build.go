package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func OnBuild() {
	OnPackage()

	build := exec.Command(Go("."), "build", "-o="+filepath.Join(".gen", "bin", "app"), ".")
	build.Env = os.Environ()

	if strings.ToLower(*FlagPlatform) == "linux/amd64" {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=amd64")
	} else if strings.ToLower(*FlagPlatform) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=arm64")
	} else if strings.ToLower(*FlagPlatform) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=amd64")
	} else if strings.ToLower(*FlagPlatform) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=arm64")
	}

	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Stdin = os.Stdin
	buildError := build.Run()
	if buildError != nil {
		Fatal(buildError)
	}
	Success("project built into ", filepath.Join(".gen", "bin", "app"))
}
