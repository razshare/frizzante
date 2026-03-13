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
	if len(ents) != 2 {
		t.Fatal("ents should contain 2 items")
	}
	if !slices.Contains(ents, "dir/subdir/test.txt") {
		t.Fatal("ents should contain dir/subdir/test.txt")
	}
	if !slices.Contains(ents, "dir/test.txt") {
		t.Fatal("ents should contain dir/test.txt")
	}
}
