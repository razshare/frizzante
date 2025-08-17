package action

import (
	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Build(o BuildOptions) error {
	err := Pkg(PkgOptions{App: o.App, Bun: o.Bun})
	if err != nil {
		return err
	}

	build := exec.Command(o.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".")
	build.Env = os.Environ()

	if o.Platform == platform.PlatformLinuxAmd64 {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=amd64")
	} else if o.Platform == platform.PlatformLinuxArm64 {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=arm64")
	} else if o.Platform == platform.PlatformDarwinAmd64 {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=amd64")
	} else if o.Platform == platform.PlatformDarwinArm64 {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=arm64")
	} else if o.Platform == platform.PlatformWindowsAmd64 {
		build.Env = append(build.Env, "GOOS=windows", "GOARCH=amd64")
	} else if o.Platform == platform.PlatformWindowsArm64 {
		build.Env = append(build.Env, "GOOS=windows", "GOARCH=arm64")
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
