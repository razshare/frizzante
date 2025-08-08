package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func OnDev() {
	OnTouch()

	mkdirError := os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm)
	if mkdirError != nil {
		Fatal(mkdirError)
	}

	air := exec.Command(Air("."))
	air.Env = append(os.Environ(), "DEV=1")
	air.Stderr = os.Stderr
	air.Stdout = os.Stdout
	air.Stdin = os.Stdin
	airError := air.Start()
	if airError != nil {
		Fatalf("air watcher faield to launch\n%s", airError)
	}
	Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { OnPackageWatch() }()

	group.Wait()
	tidyWaitError := air.Wait()
	if tidyWaitError != nil {
		Fatal(tidyWaitError)
	}
}
