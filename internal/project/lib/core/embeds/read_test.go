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

//go:embed dir/test.txt
var TestReadFileInChunksEfs embed.FS

func TestReadFileInChunks(t *testing.T) {
	var i int
	var v string
	err := ReadFileInChunks(TestReadFileInChunksEfs, "dir/test.txt", 5, func(b []byte) error {
		i++
		v += string(b)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if i != 3 {
		t.Fatal("counter should be 3")
	}

	if v != "hello world" {
		t.Fatal("value should be hello world")
	}
}
