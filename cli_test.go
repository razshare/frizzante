package main

import (
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"testing"
)

func TestOnHelp(t *testing.T) {
	cli.OnHelp()
}

func TestOnVersion(t *testing.T) {
	cli.OnVersion()
}

func TestOnCreateProject(t *testing.T) {
	directoryName := filepath.Join(".gen", "test_project")

	err := os.RemoveAll(directoryName)
	if err != nil {
		t.Fatal(err)
	}

	cli.OnCreateProject(directoryName)

	if !files.IsDirectory(directoryName) {
		t.Fatal("the cli failed to create a test project")
	}

	if !files.IsDirectory(filepath.Join(directoryName, "app")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the app directory")
	}

	if !files.IsFile(filepath.Join(directoryName, "main.go")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the main.go file")
	}

	if !files.IsFile(filepath.Join(directoryName, "go.mod")) {
		t.Fatal("the cli succeeded in creating a test project, but it's missing the go.mod file")
	}

	err = os.RemoveAll(directoryName)
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

	cli.OnAddFeature("air")

	if !files.IsDirectory(air) {
		t.Fatal("the cli failed to add air feature")
	}
}

func TestOnPackage(t *testing.T) {
	dist := filepath.Join("app", "dist")
	err := os.RemoveAll(dist)
	if err != nil {
		t.Fatal(err)
	}

	cli.OnPackage()

	if !files.IsDirectory(dist) {
		t.Fatal("the cli failed to package the application into dist")
	}
}

func TestOnInstall(t *testing.T) {
	nodeModules := filepath.Join("app", "node_modules")
	err := os.RemoveAll(nodeModules)
	if err != nil {
		t.Fatal(err)
	}

	cli.OnInstall()

	if !files.IsDirectory(nodeModules) {
		t.Fatal("the cli failed to install node_modules")
	}
}

func TestOnFormat(t *testing.T) {
	cli.OnFormat()
}

func TestOnTouch(t *testing.T) {
	dist := filepath.Join("app", "dist")
	err := os.RemoveAll(dist)
	if err != nil {
		t.Fatal(err)
	}

	cli.OnTouch()

	serverJs := filepath.Join("app", "dist", "server.js")
	if !files.IsFile(serverJs) {
		t.Fatalf("the cli failed to touch %s", serverJs)
	}

	indexHtml := filepath.Join("app", "dist", "client", "index.html")
	if !files.IsFile(indexHtml) {
		t.Fatalf("the cli failed to touch %s", indexHtml)
	}

	// We need to restore dist, otherwise other tests will break.
	cli.OnPackage()
}

func TestOnClean(t *testing.T) {
	dist := filepath.Join("app", "dist")
	nodeModules := filepath.Join("app", "node_modules")
	tmp := filepath.Join(".gen", "tmp")

	err := os.MkdirAll(dist, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(nodeModules, os.ModePerm)
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

	err = os.WriteFile(filepath.Join(nodeModules, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(tmp, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	cli.OnClean()

	if files.IsFile(filepath.Join(dist, "test.txt")) {
		t.Fatalf("cli failed to clean %s", dist)
	}

	if files.IsFile(filepath.Join(nodeModules, "test.txt")) {
		t.Fatalf("cli failed to clean %s", nodeModules)
	}

	if files.IsFile(filepath.Join(tmp, "test.txt")) {
		t.Fatalf("cli failed to clean %s", tmp)
	}

	// We need to restore dist, otherwise other tests will break.
	cli.OnInstall()
	cli.OnPackage()
}
