package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func PackageWatch() {
	Touch()

	ssr := exec.Command(path.Bun(*state.App), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=frizzante/core/scripts/server.ts")
	ssr.Dir = *state.App
	ssr.Env = append(os.Environ(), "DEV=1")
	ssr.Stderr = os.Stderr
	ssr.Stdout = os.Stdout
	ssr.Stdin = os.Stdin
	err := ssr.Start()
	if err != nil {
		messages.Fatalf("vite server watcher failed to launch\n%s", err)
	}
	messages.Success("vite server watcher launched")

	csr := exec.Command(path.Bun(*state.App), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	csr.Dir = *state.App
	csr.Env = append(os.Environ())
	csr.Stderr = os.Stderr
	csr.Stdout = os.Stdout
	csr.Stdin = os.Stdin
	err = csr.Start()
	if err != nil {
		messages.Fatalf("vite client watcher failed to launch\n%s", err)
	}
	messages.Success("vite client watcher launched")

	err = csr.Wait()
	if err != nil {
		messages.Fatal(err)
	}

	err = ssr.Wait()
	if err != nil {
		messages.Fatal(err)
	}
}
