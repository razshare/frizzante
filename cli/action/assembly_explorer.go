package action

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/textviewer"
)

func AssemblyExplorer(options AssemblyExplorerOptions) (err error) {
	var name string
	if runtime.GOOS == "windows" {
		name = filepath.Join(".gen", "bin", "app.exe")
	} else {
		name = filepath.Join(".gen", "bin", "app")
	}

	if !files.IsFile(name) && !options.Auto {
		var yes bool
		if yes, err = confirm.Sendf(true, "file %s not found. Build?", name); err != nil {
			return
		}

		if !yes {
			messages.Info("skipping build")
		} else {
			if err = Build(BuildOptions{
				App:      options.App,
				Go:       options.Go,
				Bun:      options.Bun,
				Tags:     options.Tags,
				Platform: options.Platform,
			}); err != nil {
				return
			}
		}
	} else {
		if err = Build(BuildOptions{
			App:      options.App,
			Go:       options.Go,
			Bun:      options.Bun,
			Tags:     options.Tags,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	if files.IsFile(filepath.Join(".gen", "bin", "app.s")) {
		if err = os.RemoveAll(filepath.Join(".gen", "bin", "app.s")); err != nil {
			return
		}
	}

	messages.Infof("running %s tool objdump -S %s", options.Go, name)

	var stdout *os.File
	var stderr *os.File
	var done bool

	cmd := exec.Command(options.Go, "tool", "objdump", "-S", name)
	cmd.Dir = "."
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	stdout, cmd.Stdout, _ = os.Pipe()
	stderr, cmd.Stderr, _ = os.Pipe()

	var builder strings.Builder

	go func() {
		scanner := bufio.NewScanner(stderr)
		for !done && scanner.Scan() {
			_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", messages.Prefix, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for !done && scanner.Scan() {
			builder.WriteString(scanner.Text() + "\n")
		}
	}()

	if err = cmd.Run(); err != nil {
		return
	}

	if err = textviewer.Send(fmt.Sprintf("viewing %s", filepath.Join(".gen", "bin", "app")), builder.String()); err != nil {
		return
	}

	return
}
