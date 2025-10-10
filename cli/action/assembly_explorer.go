package action

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"github.com/razshare/frizzante/tui/textviewer"
)

func AssemblyExplorer(options AssemblyExplorerOptions) (err error) {
	var asm func(source *os.File, references map[string]map[string]string, online func(line string))
	asm = func(source *os.File, references map[string]map[string]string, online func(line string)) {
		var fileName string
		var functionName string
		scanner := bufio.NewScanner(source)
		for scanner.Scan() {
			line := scanner.Text()
			online(line)
			if strings.HasPrefix(line, "TEXT") {
				parts := strings.SplitN(line, " ", 3)
				count := len(parts)
				if count > 1 {
					functionName = parts[1]
				}
				if count > 2 {
					fileName = parts[2]
				}

				functionsMap, exists := references[fileName]
				if !exists {
					functionsMap = map[string]string{}
					references[fileName] = functionsMap
				}
			} else {
				functionsMap, exists := references[fileName]
				if fileName == "" || !exists {
					continue
				}

				if _, exists = functionsMap[functionName]; !exists {
					functionsMap[functionName] = ""
				}

				functionsMap[functionName] += line + "\n"
			}
		}
	}

	var name string
	if runtime.GOOS == "windows" {
		name = filepath.Join(".gen", "bin", "app.exe")
	} else {
		name = filepath.Join(".gen", "bin", "app")
	}

	build := true

	if files.IsFile(name) {
		if build, err = confirm.Sendf(options.Auto, "file %s already exists. Rebuild?", name); err != nil {
			return
		}
	}

	tags := options.Tags

	if !build {
		messages.Info("skipping build")
	} else {
		if len(tags) == 0 {
			if tags, err = tags_.Select([]search.Choice{
				{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
				{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
				{Id: "trace", Description: "enables tracing with stack.Trace()"},
				{Id: "dry", Description: "enables dry mode"},
				{Id: "types", Description: "enables types generation"},
				{Id: "other", Description: "adds custom tags"},
			}); err != nil {
				return
			}
		}

		if err = os.RemoveAll(filepath.Join(".gen", "bin", "app.s")); err != nil {
			return
		}

		if err = Build(BuildOptions{
			App:      options.App,
			Go:       options.Go,
			Bun:      options.Bun,
			Tags:     tags,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	var references = map[string]map[string]string{}

	if files.IsFile(filepath.Join(".gen", "bin", "app.s")) {
		spin := spinner.Newf("using existing assembly code from %s", filepath.Join(".gen", "bin", "app.s"))
		go spinner.Start(spin)
		var file *os.File
		if file, err = os.Open(filepath.Join(".gen", "bin", "app.s")); err != nil {
			spinner.Stop(spin)
			return
		}
		asm(file, references, func(_ string) {})
		spinner.Stop(spin)
	} else {
		spin := spinner.Newf("generating assembly code in %s", filepath.Join(".gen", "bin", "app.s"))

		go spinner.Start(spin)
		var stdout *os.File
		var stderr *os.File

		cmd := exec.Command(options.Go, "tool", "objdump", "-S", name)
		cmd.Dir = "."
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		stdout, cmd.Stdout, _ = os.Pipe()
		stderr, cmd.Stderr, _ = os.Pipe()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", messages.Prefix, scanner.Text())
			}
		}()

		var assemblyFile *os.File
		if assemblyFile, err = os.Create(filepath.Join(".gen", "bin", "app.s")); err != nil {
			return
		}

		go asm(stdout, references, func(line string) {
			if _, werr := assemblyFile.WriteString(line + "\n"); werr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", messages.Prefix, werr.Error())
				return
			}
		})

		if err = cmd.Run(); err != nil {
			if err = assemblyFile.Close(); err != nil {
				return
			}
			spinner.Stop(spin)
			return
		}

		if err = assemblyFile.Close(); err != nil {
			return
		}
		spinner.Stop(spin)
	}

	filesChoices := make([]search.Choice, len(references))
	for {
		index := 0
		for fileName, _ := range references {
			filesChoices[index] = search.Choice{Id: fileName, Description: "explore file"}
			index++
		}

		var fileName string
		if fileName, err = singleselect.Sendf(filesChoices, "pick a file to explore"); err != nil {
			return
		}

		if fileName == "" {
			break
		}

		functionsChoices := make([]search.Choice, len(references[fileName]))

		index = 0
		for functionName, _ := range references[fileName] {
			functionsChoices[index] = search.Choice{Id: functionName, Description: "explore function"}
			index++
		}

		var functionName string
		if functionName, err = singleselect.Send(functionsChoices, "pick a function to view"); err != nil {
			return
		}

		if functionName == "" {
			break
		}

		body, exists := references[fileName][functionName]
		if !exists {
			err = fmt.Errorf("function %s not found in file %s", functionName, fileName)
			return
		}

		if err = textviewer.Send(fmt.Sprintf("viewing %s", functionName), body); err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				return
			}
		}
	}

	return
}
