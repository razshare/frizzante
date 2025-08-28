package extension

import (
	"path/filepath"
	"testing"
)

func TestFind(t *testing.T) {
	ext := Find()

	if string(filepath.Separator) == "\\" && ext != ".exe" {
		t.Fatal("extension should be .exe")
	}
}
