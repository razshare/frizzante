package apps

import (
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/cli/extensions"
)

func TestNew(t *testing.T) {
	a := New()
	if *a.Add {
		t.Fatal("add should not be active")
	}

	if *a.Help {
		t.Fatal("help should be false")
	}

	if *a.Version {
		t.Fatal("version should be false")
	}

	if *a.Reset {
		t.Fatal("reset should be false")
	}

	if *a.CreateProject {
		t.Fatal("create project should not be active")
	}

	if *a.Generate {
		t.Fatal("generate should be false")
	}

	if *a.Value != "" {
		t.Fatal("value should be empty")
	}

	if *a.Package {
		t.Fatal("package should be false")
	}

	if *a.PackageWatch {
		t.Fatal("package watch should be false")
	}

	if *a.Check {
		t.Fatal("check should be false")
	}

	if *a.Update {
		t.Fatal("update should be false")
	}

	if *a.Install {
		t.Fatal("install should be false")
	}

	if *a.Format {
		t.Fatal("format should be false")
	}

	if *a.Touch {
		t.Fatal("touch should be false")
	}

	if *a.CleanProject {
		t.Fatal("clean project should be false")
	}

	if *a.Dev {
		t.Fatal("dev should be false")
	}

	if *a.Build {
		t.Fatal("build should be false")
	}

	if *a.Configure {
		t.Fatal("configure should be false")
	}

	if *a.Platform != "" {
		t.Fatal("platform should be empty")
	}

	if *a.Yes {
		t.Fatal("yes should be false")
	}

	if *a.Go != "go"+extensions.Find() {
		t.Fatal("go should be empty")
	}

	if *a.Air != filepath.Join(".gen", "air", "air"+extensions.Find()) {
		t.Fatal("air should be empty")
	}

	if *a.Bun != filepath.Join(".gen", "bun", "bun"+extensions.Find()) {
		t.Fatal("bun should be empty")
	}

	if *a.Sqlc != filepath.Join(".gen", "sqlc", "sqlc"+extensions.Find()) {
		t.Fatal("sqlc should be empty")
	}

	if *a.Welcome {
		t.Fatal("welcome should be false")
	}

	if *a.Clear {
		t.Fatal("clear should be false")
	}
}
