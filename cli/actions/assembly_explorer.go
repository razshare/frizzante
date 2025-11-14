package actions

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/assemblies"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/hexviewer"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

func AssemblyExplorer(options AssemblyExplorerOptions) (err error) {
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
			Go:       options.Go,
			Bun:      options.Bun,
			Tags:     tags,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	var references = map[string]map[string]*assemblies.FunctionInfo{}

	if files.IsFile(filepath.Join(".gen", "bin", "app.s")) {
		spin := spinners.Newf("using existing assembly code from %s", filepath.Join(".gen", "bin", "app.s"))
		go spinners.Start(spin)
		var file *os.File
		if file, err = os.Open(filepath.Join(".gen", "bin", "app.s")); err != nil {
			spinners.Stop(spin)
			return
		}
		assemblies.ParseFunctionsInFile(file, references, func(_ string) {})
		spinners.Stop(spin)
	} else {
		spin := spinners.Newf("generating assembly code in %s", filepath.Join(".gen", "bin", "app.s"))

		go spinners.Start(spin)
		var stdout *os.File
		var stderr *os.File

		cmd := exec.Command(options.Go, "tool", "objdump", "-S", name)
		cmd.Dir = ""
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

		go assemblies.ParseFunctionsInFile(stdout, references, func(line string) {
			if _, werr := assemblyFile.WriteString(line + "\n"); werr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\r%s%s\n\r", messages.Prefix, werr.Error())
				return
			}
		})

		if err = cmd.Run(); err != nil {
			if err = assemblyFile.Close(); err != nil {
				return
			}
			spinners.Stop(spin)
			return
		}

		if err = assemblyFile.Close(); err != nil {
			return
		}
		spinners.Stop(spin)
	}

	filesChoices := make([]search.Choice, len(references))
	for {
		index := 0
		for fileName, file := range references {
			filesChoices[index] = search.Choice{Id: fileName, Description: fmt.Sprintf("%d functions", len(file))}
			index++
		}

		var fileName string
		if fileName, err = select_one.Sendf(filesChoices, "pick a file to explore"); err != nil {
			return
		}

		if fileName == "" {
			break
		}

		functionsChoices := make([]search.Choice, len(references[fileName]))

		index = 0
		for functionName, function := range references[fileName] {
			var description string

			if function.BinarySize > 1024 {
				description = fmt.Sprintf("%dKB", function.BinarySize/1024)
			} else {
				description = fmt.Sprintf("%dB", function.BinarySize)
			}

			functionsChoices[index] = search.Choice{Id: functionName, Description: description}
			index++
		}

		var functionName string
		if functionName, err = select_one.Send(functionsChoices, "pick a function to view"); err != nil {
			return
		}

		if functionName == "" {
			break
		}

		var exists bool
		var function *assemblies.FunctionInfo
		if function, exists = references[fileName][functionName]; !exists {
			err = fmt.Errorf("function %s not found in file %s", functionName, fileName)
			return
		}

		var title string

		if function.BinarySize > 1024 {
			title = fmt.Sprintf("viewing %s (%dKB)", functionName, function.BinarySize/1024)
		} else {
			title = fmt.Sprintf("viewing %s (%dB)", functionName, function.BinarySize)
		}

		formattedContent := assemblies.FormatAssemblyWithAlignment(function.AssemblyContent)

		if err = hexviewer.Send(title, formattedContent); err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				return
			}
		}
	}

	return
}
