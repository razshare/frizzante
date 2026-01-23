package types

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Generate[T any](directory string) (err error) {
	var value T
	type_ := reflect.TypeOf(value)
	var packages = map[string][]string{}
	var definitions = map[string]map[string][]string{}
	if _, err = Define(type_, packages, definitions); err != nil {
		return
	}
	if !files.IsDirectory(filepath.Join(directory)) {
		if err = os.MkdirAll(filepath.Join(directory), os.ModePerm); err != nil {
			return
		}
	}
	befores := []string{
		"github.com/razshare/frizzante/internal/project",
		"github.com/razshare/frizzante/internal/additions",
	}
	after := "main"
	packagePath := type_.PkgPath()

	for _, before := range befores {
		packagePath = strings.ReplaceAll(packagePath, before, after)
	}
	dname := filepath.Join(directory, strings.ReplaceAll(packagePath, "/", string(filepath.Separator)))
	if files.IsDirectory(dname) {
		if err = os.RemoveAll(dname); err != nil {
			return
		}
	}
	if err = os.MkdirAll(dname, os.ModePerm); err != nil {
		return
	}
	parts := strings.Split(type_.PkgPath(), "/")
	count := len(parts)
	package_ := parts[count-1]
	var globalBuilder strings.Builder
	var namespaceBuilder strings.Builder
	globalBuilder.WriteString(fmt.Sprintf("export type %s = %s.%s\n\n", type_.Name(), package_, type_.Name()))
	for namespace, definition := range definitions {
		namespaceBuilder.Reset()
		namespaceBuilder.WriteString(fmt.Sprintf("export declare namespace %s {\n", namespace))
		for name, lines := range definition {
			namespaceBuilder.WriteString(fmt.Sprintf("    export type %s = {\n", name))
			for _, line := range lines {
				namespaceBuilder.WriteString(fmt.Sprintf("        %s\n", line))
			}
			namespaceBuilder.WriteString("    }\n")
		}
		namespaceBuilder.WriteString("}\n\n")
		globalBuilder.WriteString(strings.TrimSpace(namespaceBuilder.String()))
		globalBuilder.WriteString("\n\n")
	}
	var name string
	for index, char := range type_.Name() {
		if !unicode.IsUpper(char) {
			name += string(char)
			continue
		}
		if index != 0 {
			name += "_"
		}
		name += strings.ToLower(string(char))
	}
	fname := filepath.Join(dname, name+".d.ts")
	if err = os.WriteFile(fname, []byte(strings.TrimSpace(globalBuilder.String())), os.ModePerm); err != nil {
		return
	}
	return
}
