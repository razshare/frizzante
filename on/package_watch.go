package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func PackageWatch() {
	Touch()

	server := exec.Command(path.Bun(*flags.App), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = *flags.App
	server.Env = append(os.Environ(), "DEV=1")
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Start()
	if serverError != nil {
		messages.Fatalf("vite server watcher failed to launch\n%s", serverError)
	}
	messages.Success("vite server watcher launched")

	client := exec.Command(path.Bun(*flags.App), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	client.Dir = *flags.App
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Start()
	if clientError != nil {
		messages.Fatalf("vite client watcher failed to launch\n%s", clientError)
	}
	messages.Success("vite client watcher launched")

	clientWaitError := client.Wait()
	if clientWaitError != nil {
		messages.Fatal(clientWaitError)
	}

	serverWaitError := server.Wait()
	if serverWaitError != nil {
		messages.Fatal(serverWaitError)
	}
}
