package embeds

import (
	"embed"
	"testing"
)

//go:embed is_file_test.go
var TestIsFileEfs embed.FS

func TestIsFile(t *testing.T) {
	if !IsFile(TestIsFileEfs, "is_file_test.go") {
		t.Fatal("is_file_test.go should be a file")
	}

	if IsFile(TestIsFileEfs, "qwerty") {
		t.Fatal("qwerty should not be a file")
	}
}
