package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Package() {
	Touch()

	server := exec.Command(path.Bun(*flags.App), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = *flags.App
	server.Env = append(os.Environ())
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Run()
	if serverError != nil {
		messages.Fatal(serverError)
	}

	client := exec.Command(path.Bun(*flags.App), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	client.Dir = *flags.App
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Run()
	if clientError != nil {
		messages.Fatal(clientError)
	}

	esbuild := exec.Command("node_modules/.bin/esbuild", "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esbuild.Dir = *flags.App
	esbuild.Env = append(os.Environ())
	esbuild.Stderr = os.Stderr
	esbuild.Stdout = os.Stdout
	esbuild.Stdin = os.Stdin
	esbuildError := esbuild.Run()
	if esbuildError != nil {
		messages.Fatal(esbuildError)
	}

	messages.Success("project app package generated in ", filepath.Join(*flags.App, "dist"))
}
