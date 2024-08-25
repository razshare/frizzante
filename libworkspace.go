package frizzante

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/google/uuid"
	"os"
	"regexp"
	"strings"
)
import v8 "rogchap.com/v8go"

type Workspace struct {
	NodeModulesDirectory string
	TemporaryDirectory   string
	V8Context            *v8.Context
}

func WorkspaceCreate() *Workspace {
	return &Workspace{
		NodeModulesDirectory: "",
		TemporaryDirectory:   "",
	}
}

func WorkspaceWithNodeModulesDirectory(self *Workspace, nodeModulesDirectory string) {
	self.NodeModulesDirectory = nodeModulesDirectory
}

func WorkspaceWithTemporaryDirectory(self *Workspace, temporaryDirectory string) {
	self.TemporaryDirectory = temporaryDirectory
}

type Bundle struct {
	FileName string
	Contents []byte
}

func WorkspaceBundle(self *Workspace, source []byte) (*Bundle, error) {
	dirName := self.TemporaryDirectory
	fileName := dirName + uuid.NewString() + ".js"

	mkdirAllError := os.MkdirAll(dirName, 0777)
	if mkdirAllError != nil {
		return nil, mkdirAllError
	}

	err := os.WriteFile(fileName, source, 0777)

	if err != nil {
		return nil, err
	}

	buildResult := api.Build(api.BuildOptions{
		Bundle:      true,
		Format:      api.FormatESModule,
		EntryPoints: []string{fileName},
		NodePaths:   []string{self.NodeModulesDirectory},
	})

	for _, buildError := range buildResult.Errors {
		builder := strings.Builder{}

		builder.WriteString(buildError.Text)
		builder.WriteString("\n")
		builder.WriteString(buildError.Location.LineText)
		builder.WriteString("\n")
		builder.WriteString("(")
		builder.WriteString(string(rune(buildError.Location.Line)))
		builder.WriteString(":")
		builder.WriteString(string(rune(buildError.Location.Column)))
		builder.WriteString(")")
		builder.WriteString("\n")
		builder.WriteString(buildError.Location.File)

		return nil, fmt.Errorf(builder.String())
	}

	for _, file := range buildResult.OutputFiles {
		removeError := os.Remove(fileName)
		if removeError != nil {
			return nil, removeError
		}

		stringifiedContents := string(file.Contents)
		replaced := strings.Replace(stringifiedContents, "\"use strict\";", "", 1)
		concat := requirements + replaced

		return &Bundle{Contents: []byte(concat), FileName: fileName}, nil
	}

	removeError := os.Remove(fileName)
	if removeError != nil {
		return nil, removeError
	}

	return nil, fmt.Errorf("could not build input file")
}

var regexSsr, regexSsrError = regexp.Compile(`var \w+[\s\n]*=[\s\n]*Component[\s\n]*;?[\s\n]*export[\s\n]+{[\s\n]*\w+[\s\n]+as[\s\n]+default[\s\n]*}[\s\n]*;?`)

func WorkspaceCompileSvelte(self *Workspace, svelteFileName string) (func(props map[string]any) (string, error), error) {
	if regexSsrError != nil {
		return nil, regexSsrError
	}

	bundle, bundleError := WorkspaceBundle(self, boot)
	if bundleError != nil {
		return nil, bundleError
	}
	isoLocal := v8.NewIsolate()

	externGetArgs := v8.NewFunctionTemplate(isoLocal, func(info *v8.FunctionCallbackInfo) *v8.Value {
		indexContents, indexError := os.ReadFile(svelteFileName)
		if indexError != nil {
			return nil
		}

		value, valueError := v8.NewValue(isoLocal, string(indexContents))
		if valueError != nil {
			return nil
		}

		return value
	})

	globalLocal := v8.NewObjectTemplate(isoLocal)
	externGetArgsError := globalLocal.Set("externGetArgs", externGetArgs)
	if externGetArgsError != nil {
		return nil, externGetArgsError
	}

	contextLocal := v8.NewContext(isoLocal, globalLocal)

	script := string(bundle.Contents)

	compileResult, runError := contextLocal.RunScript(script, "compile.js")
	if runError != nil {
		return nil, runError
	}

	compiledScript := compileResult.String()

	ssrBundle, renderBundleError := WorkspaceBundle(self, []byte(compiledScript))
	if renderBundleError != nil {
		return nil, renderBundleError
	}

	ssrScript := string(regexSsr.ReplaceAll(ssrBundle.Contents, []byte("Component.render")))

	render, ssrScriptError := contextGlobal.RunScript(ssrScript, svelteFileName)
	if ssrScriptError != nil {
		return nil, ssrScriptError
	}

	function, functionError := render.AsFunction()

	if functionError != nil {
		return nil, functionError
	}

	return func(props map[string]any) (string, error) {
		objectTemplate := v8.NewObjectTemplate(isolateGlobal)
		for key, value := range props {
			objectTemplate.Set(key, value)
		}

		instance, instanceError := objectTemplate.NewInstance(contextGlobal)
		if instanceError != nil {
			return "", instanceError
		}
		returnValue, callError := function.Call(contextGlobal.Global(), instance)
		if callError != nil {
			return "", callError
		}

		output, outputError := returnValue.AsObject()

		if outputError != nil {
			return "", outputError
		}

		htmlValue, htmlValueError := output.Get("html")

		if htmlValueError != nil {
			return "", htmlValueError
		}

		html := htmlValue.String()

		return html, nil
	}, nil
}
