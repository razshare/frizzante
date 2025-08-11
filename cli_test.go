package main

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/on"
	"os"
	"path/filepath"
	"testing"
)

func TestOnHelp(test *testing.T) {
	on.Help()
}

func TestOnVersion(test *testing.T) {
	on.Version(efs)
}

func TestOnCreateProject(test *testing.T) {
	directoryName := filepath.Join(".gen", "test_project")

	err := os.RemoveAll(directoryName)
	if err != nil {
		test.Fatal(err)
	}

	on.CreateProject(directoryName)

	if !files.IsDirectory(directoryName) {
		test.Fatal("the cli failed to create a test project")
	}

	if !files.IsDirectory(filepath.Join(directoryName, "app")) {
		test.Fatal("the cli succeeded in creating a test project, but it's missing the app directory")
	}

	if !files.IsFile(filepath.Join(directoryName, "main.go")) {
		test.Fatal("the cli succeeded in creating a test project, but it's missing the main.go file")
	}

	if !files.IsFile(filepath.Join(directoryName, "go.mod")) {
		test.Fatal("the cli succeeded in creating a test project, but it's missing the go.mod file")
	}

	err = os.RemoveAll(directoryName)
	if err != nil {
		test.Fatal(err)
	}
}

func TestOnAddFeature(test *testing.T) {
	air := filepath.Join(".gen", "air")
	err := os.RemoveAll(air)
	if err != nil {
		test.Fatal(err)
	}

	on.Generate(efs, "air")

	if !files.IsDirectory(air) {
		test.Fatal("the cli failed to add air feature")
	}
}

func TestOnPackage(test *testing.T) {
	dist := filepath.Join("template", "app", "dist")
	err := os.RemoveAll(dist)
	if err != nil {
		test.Fatal(err)
	}

	on.Package()

	if !files.IsDirectory(dist) {
		test.Fatal("the cli failed to package the application into dist")
	}
}

func TestOnInstall(test *testing.T) {
	nodeModules := filepath.Join("template", "app", "node_modules")
	err := os.RemoveAll(nodeModules)
	if err != nil {
		test.Fatal(err)
	}

	on.Install()

	if !files.IsDirectory(nodeModules) {
		test.Fatal("the cli failed to install node_modules")
	}
}

func TestOnFormat(test *testing.T) {
	on.Format()
}

func TestOnTouch(test *testing.T) {
	dist := filepath.Join("template", "app", "dist")
	err := os.RemoveAll(dist)
	if err != nil {
		test.Fatal(err)
	}

	on.Touch()

	serverJs := filepath.Join("template", "app", "dist", "server.js")
	if !files.IsFile(serverJs) {
		test.Fatalf("the cli failed to touch %s", serverJs)
	}

	indexHtml := filepath.Join("template", "app", "dist", "client", "index.html")
	if !files.IsFile(indexHtml) {
		test.Fatalf("the cli failed to touch %s", indexHtml)
	}

	// We need to restore dist, otherwise other tests will break.
	on.Package()
}

func TestOnClean(test *testing.T) {
	dist := filepath.Join("template", "app", "dist")
	nodeModules := filepath.Join("template", "app", "node_modules")
	tmp := filepath.Join(".gen", "tmp")

	err := os.MkdirAll(dist, os.ModePerm)
	if err != nil {
		test.Fatal(err)
	}

	err = os.MkdirAll(nodeModules, os.ModePerm)
	if err != nil {
		test.Fatal(err)
	}

	err = os.MkdirAll(tmp, os.ModePerm)
	if err != nil {
		test.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(dist, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		test.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(nodeModules, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		test.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(tmp, "test.txt"), []byte("test"), os.ModePerm)
	if err != nil {
		test.Fatal(err)
	}

	on.Clean()

	if files.IsFile(filepath.Join(dist, "test.txt")) {
		test.Fatalf("cli failed to clean %s", dist)
	}

	if files.IsFile(filepath.Join(nodeModules, "test.txt")) {
		test.Fatalf("cli failed to clean %s", nodeModules)
	}

	if files.IsFile(filepath.Join(tmp, "test.txt")) {
		test.Fatalf("cli failed to clean %s", tmp)
	}

	// We need to restore dist, otherwise other tests will break.
	on.Install()
	on.Package()
}
