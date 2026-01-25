package embeds

import (
	"embed"
	"testing"
)

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
