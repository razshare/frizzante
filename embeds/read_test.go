package embeds

import (
	"embed"
	"slices"
	"testing"
)

//go:embed dir
var TestReadDirectoryEfs embed.FS

func TestReadDirectory(t *testing.T) {
	ents, err := ReadDirectory(TestReadDirectoryEfs, "dir")
	if err != nil {
		t.Fatal(err)
	}

	if len(ents) != 1 {
		t.Fatal("dir should contain only 1 file")
	}

	if !slices.Contains(ents, "dir/test.txt") {
		t.Fatal("dir/test.txt")
	}
}
