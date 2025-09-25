package messages

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func Command(dir string, env []string, name string, args ...string) (ok bool) {
	var stdout *os.File
	var stderr *os.File
	var done bool
	defer func() { done = true }()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = os.Stdin

	stdout, cmd.Stdout, _ = os.Pipe()
	stderr, cmd.Stderr, _ = os.Pipe()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for !done && scanner.Scan() {
			_, _ = fmt.Fprintf(os.Stdout, "\r%s%s\n\r", Prefix, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for !done && scanner.Scan() {
			_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", Prefix, scanner.Text())
		}
	}()

	if err := cmd.Run(); err != nil {
		ok = false
		return
	}

	ok = true

	return
}
