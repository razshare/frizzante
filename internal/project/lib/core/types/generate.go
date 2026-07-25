package types

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
)

func Generate[T any]() (err error) {
	typesDirectoryName := filepath.Join("app", "lib", "types", "server")
	var value T
	type_ := reflect.TypeOf(value)
	var packages = map[string][]string{}
	var definitions = map[string]map[string][]string{}
	if _, err = Define(type_, packages, definitions); err != nil {
		return
	}
	if !files.IsDirectory(filepath.Join(typesDirectoryName)) {
		if err = os.MkdirAll(filepath.Join(typesDirectoryName), os.ModePerm); err != nil {
			return
		}
	}
	befores := []string{
		"main",
	}
	after := "main"
	packagePath := type_.PkgPath()
	for _, before := range befores {
		packagePath = strings.ReplaceAll(packagePath, before, after)
	}
	packageDirectoryName := filepath.Join(typesDirectoryName, strings.ReplaceAll(packagePath, "/", string(filepath.Separator)))
	if !files.IsDirectory(packageDirectoryName) {
		if err = os.MkdirAll(packageDirectoryName, os.ModePerm); err != nil {
			return
		}
	}
	parts := strings.Split(type_.PkgPath(), "/")
	count := len(parts)
	package_ := parts[count-1]
	var globalBuilder strings.Builder
	var namespaceBuilder strings.Builder
	if _, err = fmt.Fprintf(&globalBuilder, "export type %s = %s.%s\n\n", type_.Name(), package_, type_.Name()); err != nil {
		return
	}
	for namespace, definition := range definitions {
		namespaceBuilder.Reset()
		if _, err = fmt.Fprintf(&namespaceBuilder, "export declare namespace %s {\n", namespace); err != nil {
			return
		}
		for name, lines := range definition {
			if _, err = fmt.Fprintf(&namespaceBuilder, "    export type %s = {\n", name); err != nil {
				return
			}
			for _, line := range lines {
				if _, err = fmt.Fprintf(&namespaceBuilder, "        %s\n", line); err != nil {
					return
				}
			}
			namespaceBuilder.WriteString("    }\n")
		}
		namespaceBuilder.WriteString("}\n\n")
		globalBuilder.WriteString(strings.TrimSpace(namespaceBuilder.String()))
		globalBuilder.WriteString("\n\n")
	}
	var name strings.Builder
	for index, char := range type_.Name() {
		if !unicode.IsUpper(char) {
			name.WriteString(string(char))
			continue
		}
		if index != 0 {
			name.WriteString("_")
		}
		name.WriteString(strings.ToLower(string(char)))
	}
	fileName := filepath.Join(packageDirectoryName, name.String()+".d.ts")
	if err = os.WriteFile(fileName, []byte(strings.TrimSpace(globalBuilder.String())), os.ModePerm); err != nil {
		return
	}
	return
}
