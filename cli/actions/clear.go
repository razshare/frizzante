package actions

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Clear(_ ClearOptions) (err error) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		return
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		return
	}

	if runtime.GOOS == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		return
	}

	fmt.Print("\033[H\033[2J")
	return
}
