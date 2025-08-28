package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMove(t *testing.T) {
	_ = os.RemoveAll(filepath.Join("dir", "move"))
	defer func() { _ = os.RemoveAll(filepath.Join("dir", "move")) }()

	err := os.WriteFile(filepath.Join("dir", "move.txt"), make([]byte, 0), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = Move(filepath.Join("dir", "move.txt"), filepath.Join("dir", "move", "move.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if IsFile(filepath.Join("dir", "move.txt")) {
		t.Fatal("dir/move.txt should not be a file")
	}

	if !IsDirectory(filepath.Join("dir", "move")) {
		t.Fatal("dir/move should be a directory")
	}
}

func TestMoveShouldFail(t *testing.T) {
	if err := Move("-- ;,", filepath.Join("dir", "move", "move.txt")); err == nil {
		t.Fatal("move should fail")
	}
}
