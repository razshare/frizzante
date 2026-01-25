package extensions

import (
	"path/filepath"
	"testing"
)

func TestFind(t *testing.T) {
	if ext := Find(); string(filepath.Separator) == "\\" && ext != ".exe" {
		t.Fatal("extension should be .exe")
	}
}
