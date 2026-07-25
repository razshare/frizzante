package actions

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func LockPackages(_ LockPackagesOptions) (err error) {
	type PackageJson struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if !files.IsFile(filepath.Join("app", "package.json")) {
		err = errors.New("app/package.json not found")
		return
	}
	var data []byte
	var packageJson PackageJson
	if data, err = os.ReadFile(filepath.Join("app", "package.json")); err != nil {
		return err
	}
	if err = json.Unmarshal(data, &packageJson); err != nil {
		return
	}
	var devDependenciesFixed map[string]string
	if packageJson.DevDependencies != nil {
		devDependenciesFixed = map[string]string{}
		for packageName, packageVersionString := range packageJson.DevDependencies {
			packageVersionString = strings.ReplaceAll(packageVersionString, "^", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, ">=", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, "<=", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, "<", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, ">", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, "~", "")
			devDependenciesFixed[packageName] = packageVersionString
		}
	}
	var dependenciesFixed map[string]string
	if packageJson.Dependencies != nil {
		dependenciesFixed = map[string]string{}
		for packageName, packageVersionString := range packageJson.Dependencies {
			packageVersionString = strings.ReplaceAll(packageVersionString, "^", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, ">=", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, "<=", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, "<", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, ">", "")
			packageVersionString = strings.ReplaceAll(packageVersionString, "~", "")
			dependenciesFixed[packageName] = packageVersionString
		}
	}
	var packageJsonUnknown map[string]any
	if err = json.Unmarshal(data, &packageJsonUnknown); err != nil {
		return
	}
	if _, exists := packageJsonUnknown["dependencies"]; exists && dependenciesFixed != nil {
		packageJsonUnknown["dependencies"] = dependenciesFixed
	}
	if _, exists := packageJsonUnknown["devDependencies"]; exists && devDependenciesFixed != nil {
		packageJsonUnknown["devDependencies"] = devDependenciesFixed
	}
	if data, err = json.MarshalIndent(packageJsonUnknown, "", "    "); err != nil {
		return
	}
	err = os.WriteFile(filepath.Join("app", "package.json"), data, os.ModePerm)
	messages.Success("packages locked")
	return
}
