package messages

import (
	"os"
	"os/exec"
	"strings"
	"sync"
)

var Mutex sync.Mutex

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

	transfer := func(from *os.File, to *os.File, buffer []byte) bool {
		Mutex.Lock()
		defer Mutex.Unlock()

		if _, rerr := from.Read(buffer); rerr != nil {
			return false
		}

		content := string(buffer)

		for _, line := range strings.Split(content, "\n") {
			if _, werr := to.WriteString("\r" + Prefix + line + "\n"); werr != nil {
				return false
			}
		}

		return true
	}

	go func() {
		buffer := make([]byte, 1024)
		for !done && transfer(stdout, os.Stdout, buffer) {
		}
	}()

	go func() {
		buffer := make([]byte, 1024)
		for !done && transfer(stderr, os.Stderr, buffer) {
		}
	}()

	if err := cmd.Run(); err != nil {
		ok = false
		return
	}

	ok = true

	return
}
