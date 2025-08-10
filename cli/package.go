package cli

import (
	"os"
	"os/exec"
)

func OnPackage() {
	OnTouch()

	server := exec.Command(Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = "app"
	server.Env = append(os.Environ())
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Run()
	if serverError != nil {
		Fatal(serverError)
	}

	client := exec.Command(Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	client.Dir = "app"
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Run()
	if clientError != nil {
		Fatal(clientError)
	}

	esbuild := exec.Command("node_modules/.bin/esbuild", "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esbuild.Dir = "app"
	esbuild.Env = append(os.Environ())
	esbuild.Stderr = os.Stderr
	esbuild.Stdout = os.Stdout
	esbuild.Stdin = os.Stdin
	esbuildError := esbuild.Run()
	if esbuildError != nil {
		Fatal(esbuildError)
	}

	Success("project app package generated in app/dist")
}

func OnPackageWatch() {
	OnTouch()

	server := exec.Command(Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = "app"
	server.Env = append(os.Environ(), "DEV=1")
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Start()
	if serverError != nil {
		Fatalf("vite server watcher failed to launch\n%s", serverError)
	}
	Success("vite server watcher launched")

	client := exec.Command(Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	client.Dir = "app"
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Start()
	if clientError != nil {
		Fatalf("vite client watcher failed to launch\n%s", clientError)
	}
	Success("vite client watcher launched")

	clientWaitError := client.Wait()
	if clientWaitError != nil {
		Fatal(clientWaitError)
	}

	serverWaitError := server.Wait()
	if serverWaitError != nil {
		Fatal(serverWaitError)
	}
}
