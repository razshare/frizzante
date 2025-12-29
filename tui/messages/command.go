package messages

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Command(options CommandOptions) (ok bool) {
	var err error
	var program string
	if files.IsFile(options.Program) {
		if program, err = filepath.Abs(options.Program); err != nil {
			Error(err)
			return
		}
	} else if program, err = exec.LookPath(options.Program); err != nil {
		Error(err)
		return
	}

	var done bool
	defer func() { done = true }()

	cmd := exec.Command(program, options.Args...)
	cmd.Dir = options.DirectoryName
	cmd.Env = options.Environment

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

	if err = cmd.Run(); err != nil {
		Error(err)
		ok = false
		return
	}

	ok = true

	return
}
