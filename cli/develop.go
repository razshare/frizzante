package cli

import (
	"embed"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func Develop(efs embed.FS) {
	if !*FlagDevelop {
		return
	}

	if !fs.FileExists("main.go") {
		log.Fatal("no main.go file detected")
	}

	if !fs.FileExists(".air.toml") {
		data, readError := efs.ReadFile(".air.toml")
		if readError != nil {
			log.Fatal(readError)
		}

		writeError := os.WriteFile(".air.toml", data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}
	}

	if !fs.FileExists(filepath.Join("bin", "bun")) {
		// Configure bun...
		spinner, _ := pterm.DefaultSpinner.
			WithShowTimer(false).
			WithRemoveWhenDone(true).
			Start("Configuring Bun...")

		data, readError := efs.ReadFile("bin/bun")
		if readError != nil {
			log.Fatal(readError)
		}

		writeError := os.WriteFile(filepath.Join("bin", "bun"), data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}

		stopError := spinner.Stop()
		if stopError != nil {
			log.Fatal(stopError)
		}
	}

	if !fs.FileExists(filepath.Join("bin", "air")) {
		// Configure air...
		spinner, _ := pterm.DefaultSpinner.
			WithShowTimer(false).
			WithRemoveWhenDone(true).
			Start("Configuring Air...")

		data, readError := efs.ReadFile("bin/air")
		if readError != nil {
			log.Fatal(readError)
		}

		writeError := os.WriteFile(filepath.Join("bin", "air"), data, os.ModePerm)
		if writeError != nil {
			log.Fatal(writeError)
		}

		stopError := spinner.Stop()
		if stopError != nil {
			log.Fatal(stopError)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	directory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Create commands...
	///// bun update
	bunUpdate := exec.Command(filepath.Join("..", "bin", "bun"), "update")
	bunUpdate.Dir = filepath.Join(directory, "app")
	bunUpdate.Stdout = os.Stdout
	bunUpdate.Stderr = os.Stderr
	bunUpdate.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// tidy
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = filepath.Join(directory, "app")
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	tidy.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// csr
	csr := exec.Command(
		filepath.Join("..", "bin", "bun"),
		"x",
		"--bun",
		"vite",
		"build",
		"--outDir=dist/client",
		"--logLevel=info",
	)
	csr.Dir = filepath.Join(directory, "app")
	csr.Stdout = os.Stdout
	csr.Stderr = os.Stderr
	csr.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// ssr
	ssr := exec.Command(
		filepath.Join("..", "bin", "bun"),
		"x",
		"--bun",
		"vite",
		"build",
		"--watch",
		"--outDir=dist",
		"--logLevel=info",
		"--ssr=lib/utilities/frz/scripts/server.ts",
	)
	ssr.Dir = filepath.Join(directory, "app")
	ssr.Stdout = os.Stdout
	ssr.Stderr = os.Stderr
	ssr.Env = append(
		os.Environ(),
		"DEV=1",
	)
	///// air
	air := exec.Command(filepath.Join("bin", "air"))
	air.Dir = directory
	air.Stdout = os.Stdout
	air.Stderr = os.Stderr
	air.Env = append(
		os.Environ(),
		"DEV=1",
		"CGO_ENABLED=1",
	)

	// Run bun update once...
	bunUpdateError := bunUpdate.Run()
	if bunUpdateError != nil {
		log.Fatal(bunUpdateError)
	}

	// Run go mod tidy once...
	tidyError := tidy.Run()
	if tidyError != nil {
		log.Fatal(tidyError)
	}

	// Run vite csr build once...
	csrError := csr.Run()
	if csrError != nil {
		log.Fatal(csrError)
	}

	go func() {
		// Start vite ssr...
		ssrError := ssr.Start()
		if ssrError != nil {
			_ = air.Cancel()
			log.Fatal(ssrError)
		}
	}()

	go func() {
		// Start air...
		airError := air.Start()
		if airError != nil {
			_ = ssr.Cancel()
			log.Fatal(airError)
		}
	}()

	<-sigs

	_ = ssr.Wait()
	_ = air.Wait()

	os.Exit(0)
}
