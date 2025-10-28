package messages

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func Command(options CommandOptions) (ok bool) {
	var done bool
	defer func() { done = true }()

	cmd := exec.Command(options.Name, options.Args...)
	cmd.Dir = options.Dir
	cmd.Env = options.Env

	if !options.DisabledStdin {
		cmd.Stdin = os.Stdin
	}

	if !options.DisableStdout {
		var stdout *os.File
		stdout, cmd.Stdout, _ = os.Pipe()
		go func() {
			scanner := bufio.NewScanner(stdout)
			for !done && scanner.Scan() {
				_, _ = fmt.Fprintf(os.Stdout, "\r%s%s\n\r", Prefix, scanner.Text())
			}
		}()
	}

	if !options.DisableStderr {
		var stderr *os.File
		stderr, cmd.Stderr, _ = os.Pipe()
		go func() {
			scanner := bufio.NewScanner(stderr)
			for !done && scanner.Scan() {
				_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", Prefix, scanner.Text())
			}
		}()
	}

	if err := cmd.Run(); err != nil {
		Error(err)
		ok = false
		return
	}

	ok = true

	return
}
