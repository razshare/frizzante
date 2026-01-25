package files

import (
	"path/filepath"
	"testing"
)

func TestReadFileInChunks(t *testing.T) {
	var index int
	var value string
	err := ReadFileInChunks(filepath.Join("dir", "test.txt"), 5, func(b []byte) error {
		index++
		value += string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if index != 3 {
		t.Fatal("counter should be 3")
	}
	if value != "hello world" {
		t.Fatal("value should be hello world")
	}
}
