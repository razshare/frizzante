package main

import (
	on2 "github.com/razshare/frizzante/cli/on"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"testing"
)

func TestOnHelp(t *testing.T) {
	err := on2.Help(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnVersion(t *testing.T) {
	err := on2.Version(c)
	if err != nil {
		t.Fatal()
	}
}

func TestOnCreateProject(t *testing.T) {
	dst := filepath.Join(".gen", "test_project")

	err := os.RemoveAll(dst)
	if err != nil {
		t.Fatal(err)
	}

	app := *c.Flags.App
	*c.Flags.App = "app"
	defer func() { *c.Flags.App = app }()
	err = on2.CreateProject(c, dst)
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(dst) {
		t.Fatal("the cli failed to create a test project")
	}

	if !files.IsDirectory(filepath.Join(dst, "app")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the app directory")
	}

	if !files.IsFile(filepath.Join(dst, "main.go")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the main.go file")
	}

	if !files.IsFile(filepath.Join(dst, "go.mod")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the go.mod file")
	}

	err = os.RemoveAll(dst)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnAddFeature(t *testing.T) {
	air := filepath.Join(".gen", "air")
	err := os.RemoveAll(air)
	if err != nil {
		t.Fatal(err)
	}

	err = on2.Generate(c, ".", "air")
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(air) {
		t.Fatal("the cli failed to add air feature")
	}
}

func TestOnPackage(t *testing.T) {
	dist := filepath.Join("template", "app", "dist")
	err := os.RemoveAll(dist)
	if err != nil {
		t.Fatal(err)
	}

	err = on2.Package(c)
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(dist) {
		t.Fatal("the cli failed to package the application into dist")
	}
}

func TestOnInstall(t *testing.T) {
	mods := filepath.Join("template", "app", "node_modules")
	err := os.RemoveAll(mods)
	if err != nil {
		t.Fatal(err)
	}

	err = on2.Install(c, ".")
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(mods) {
		t.Fatal("the cli failed to install node_modules")
	}
}

func TestOnFormat(t *testing.T) {
	err := on2.Format(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOnTouch(t *testing.T) {
	dist := filepath.Join("template", "app", "dist")
	err := os.RemoveAll(dist)
	if err != nil {
		t.Fatal(err)
	}

	err = on2.Touch(c)
	if err != nil {
		t.Fatal(err)
	}

	script := filepath.Join("template", "app", "dist", "server.js")
	if !files.IsFile(script) {
		t.Fatalf("the cli failed to touch %s", script)
	}

	doc := filepath.Join("template", "app", "dist", "client", "index.html")
	if !files.IsFile(doc) {
		t.Fatalf("the cli failed to touch %s", doc)
	}

	// We need to restore dist, otherwise other tests will break.
	err = on2.Package(c)
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

	err = on2.Clean(c)
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
	err = on2.Install(c, ".")
	if err != nil {
		t.Fatal(err)
	}
	err = on2.Package(c)
	if err != nil {
		t.Fatal(err)
	}
}
