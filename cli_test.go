package main

import (
	"github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"os"
	"path/filepath"
	"testing"
)

func TestOnHelp(t *testing.T) {
	err := action.Help(action.HelpOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnVersion(t *testing.T) {
	err := action.Version(action.VersionOptions{Efs: efs})
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnCreateProject(t *testing.T) {
	proj := filepath.Join(".gen", "test_project")
	err := os.RemoveAll(proj)
	if err != nil {
		t.Fatal(err)
	}

	err = action.CreateProject(action.CreateProjectOptions{Project: proj})
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(proj) {
		t.Fatal("the cli failed to create a test project")
	}

	if !files.IsDirectory(filepath.Join(proj, "app")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the app directory")
	}

	if !files.IsFile(filepath.Join(proj, "main.go")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the main.go file")
	}

	if !files.IsFile(filepath.Join(proj, "go.mod")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the go.mod file")
	}

	err = os.RemoveAll(proj)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnGenerate(t *testing.T) {
	err := os.RemoveAll(filepath.Join(".gen", "air"))
	if err != nil {
		t.Fatal(err)
	}

	err = action.Generate(action.GenerateOptions{
		App:      filepath.Join("template", "app"),
		Selected: "air",
		Platform: platform.PlatformLinuxAmd64,
		Auto:     true,
		Go:       "go",
		Air:      filepath.Join(".gen", "air", "air"),
		Bun:      filepath.Join(".gen", "bun", "bun"),
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
		Efs:      tefs,
	})

	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "air", "air")) {
		t.Fatal("the cli failed to add air feature")
	}
}

func TestOnPackage(t *testing.T) {
	err := os.RemoveAll(filepath.Join("template", "app", "dist"))
	if err != nil {
		t.Fatal(err)
	}

	err = action.Pkg(action.PkgOptions{
		App: filepath.Join("template", "app"),
		Bun: filepath.Join(".gen", "bun", "bun"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(filepath.Join("template", "app", "dist")) {
		t.Fatal("the cli failed to package the application into dist")
	}
}

func TestOnInstall(t *testing.T) {
	err := os.RemoveAll(filepath.Join("template", "app", "node_modules"))
	if err != nil {
		t.Fatal(err)
	}

	err = action.Install(action.InstallOptions{
		App: filepath.Join("template", "app"),
		Go:  "go",
		Bun: filepath.Join(".gen", "bun", "bun"),
	})

	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(filepath.Join("template", "app", "node_modules")) {
		t.Fatal("the cli failed to install node_modules")
	}
}

func TestOnFormat(t *testing.T) {
	err := action.Format(action.FormatOptions{
		App: filepath.Join("template", "app"),
		Go:  "go",
		Bun: filepath.Join(".gen", "bun", "bun"),
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestOnTouch(t *testing.T) {
	err := os.RemoveAll(filepath.Join("template", "app", "dist"))
	if err != nil {
		t.Fatal(err)
	}

	err = action.Touch(action.TouchOptions{App: filepath.Join("template", "app")})
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("template", "app", "dist", "server.js")) {
		t.Fatalf("the cli failed to touch %s", filepath.Join("template", "app", "dist", "server.js"))
	}

	if !files.IsFile(filepath.Join("template", "app", "dist", "client", "index.html")) {
		t.Fatalf("the cli failed to touch %s", filepath.Join("template", "app", "dist", "client", "index.html"))
	}

	// We need to restore the package, otherwise other tests will break.
	err = action.Pkg(action.PkgOptions{
		App: filepath.Join("template", "app"),
		Bun: filepath.Join(".gen", "bun", "bun"),
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestOnClean(t *testing.T) {
	dist := filepath.Join("template", "app", "dist")
	mods := filepath.Join("template", "app", "node_modules")
	tmp := filepath.Join(".gen", "tmp")

	err := os.MkdirAll(dist, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(mods, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(tmp, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(dist, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(mods, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(tmp, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = action.CleanProject(action.CleanProjectOptions{
		App: filepath.Join("template", "app"),
		Go:  "go",
	})

	if err != nil {
		t.Fatal(err)
	}

	if files.IsFile(filepath.Join(dist, "test.txt")) {
		t.Fatalf("cli failed to clean %s", dist)
	}

	if files.IsFile(filepath.Join(mods, "test.txt")) {
		t.Fatalf("cli failed to clean %s", mods)
	}

	if files.IsFile(filepath.Join(tmp, "test.txt")) {
		t.Fatalf("cli failed to clean %s", tmp)
	}

	// We need to restore dist, otherwise other tests will break.
	err = action.Install(action.InstallOptions{
		App: filepath.Join("template", "app"),
		Go:  "go",
		Bun: filepath.Join(".gen", "bun", "bun"),
	})

	if err != nil {
		t.Fatal(err)
	}

	err = action.Pkg(action.PkgOptions{
		App: filepath.Join("template", "app"),
		Bun: filepath.Join(".gen", "bun", "bun"),
	})

	if err != nil {
		t.Fatal(err)
	}
}
