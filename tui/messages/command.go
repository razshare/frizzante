package messages

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/values"
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
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, program, options.Args...)
	cmd.Dir = options.DirectoryName
	cmd.Env = options.Environment
	if !options.DisabledStdin {
		cmd.Stdin = os.Stdin
	}
	if !options.DisableStdout || options.StdoutBuilder != nil {
		var stdout *os.File
		stdout, cmd.Stdout, _ = os.Pipe()
		go func() {
			scanner := bufio.NewScanner(stdout)
			for !done && scanner.Scan() {
				if options.StdoutBuilder != nil {
					options.StdoutBuilder.WriteString(fmt.Sprintf("%s\n", scanner.Text()))
				}
				if !options.DisableStdout {
					_, _ = fmt.Fprintf(os.Stdout, "\r%s%s\n\r", Prefix, scanner.Text())
				}
			}
		}()
	}
	if !options.DisableStderr || options.StderrBuilder != nil {
		var stderr *os.File
		stderr, cmd.Stderr, _ = os.Pipe()
		go func() {
			scanner := bufio.NewScanner(stderr)
			for !done && scanner.Scan() {
				if options.StderrBuilder != nil {
					options.StderrBuilder.WriteString(fmt.Sprintf("%s\n", scanner.Text()))
				}
				if !options.DisableStderr {
					_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", Prefix, scanner.Text())
				}
			}
		}()
	}
	if cmd.Cancel != nil && options.Channels.End != nil {
		go func() {
			<-options.Channels.End
			options.Channels.End <- values.None
			if cerr := cmd.Cancel(); cerr != nil {
				Error(err)
				ok = false
				return
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
