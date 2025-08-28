package npm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstall(t *testing.T) {
	t.Run("empty packages returns nil", func(t *testing.T) {
		if err := Install("bun", "/tmp/test"); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("invalid directory returns error", func(t *testing.T) {
		err := Install("bun", "/non/existent/directory", "package1")
		if err == nil {
			t.Fatal("expected error for invalid directory")
		}
		expectedMsg := "directory /non/existent/directory not found"
		if err.Error() != expectedMsg {
			t.Fatalf("expected error message '%s', got '%v'", expectedMsg, err)
		}
	})

	t.Run("valid directory with packages", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "npm-test")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.RemoveAll(tmpDir) }()

		packageJsonPath := filepath.Join(tmpDir, "package.json")

		if err = os.WriteFile(packageJsonPath, []byte(`{"name":"test","version":"1.0.0"}`), 0644); err != nil {
			t.Fatal(err)
		}

		if err = Install("echo", tmpDir, "test-package"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		nodeModulesPath := filepath.Join(tmpDir, "node_modules")

		if _, err = os.Stat(nodeModulesPath); err == nil {
			t.Logf("node_modules created at %s - will be cleaned up", nodeModulesPath)
		}
	})

	t.Run("multiple packages", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "npm-test-multi")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.RemoveAll(tmpDir) }()

		packageJsonPath := filepath.Join(tmpDir, "package.json")

		if err = os.WriteFile(packageJsonPath, []byte(`{"name":"test","version":"1.0.0"}`), 0644); err != nil {
			t.Fatal(err)
		}

		err = Install("echo", tmpDir, "package1", "package2", "package3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		nodeModulesPath := filepath.Join(tmpDir, "node_modules")
		if _, err = os.Stat(nodeModulesPath); err == nil {
			t.Logf("node_modules created at %s - will be cleaned up", nodeModulesPath)
		}
	})

	t.Run("handles command failure gracefully", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "npm-test-fail")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = os.RemoveAll(tmpDir) }()

		if err = Install("/non/existent/command", tmpDir, "package1"); err != nil {
			t.Fatalf("Install should not return error even if command fails: %v", err)
		}

		nodeModulesPath := filepath.Join(tmpDir, "node_modules")
		if _, err = os.Stat(nodeModulesPath); err == nil {
			t.Logf("node_modules created at %s - will be cleaned up", nodeModulesPath)
		}
	})
}
