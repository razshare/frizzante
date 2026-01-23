package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/internal/project/lib/core/snapshots"
)

var directory = filepath.Join(".gen", "snapshot")

func main() {
	if err := os.RemoveAll(directory); err != nil {
		log.Fatal(err)
	}
	if err := snapshots.Generate("http://127.0.0.1:8080/", directory); err != nil {
		log.Fatal(err)
	}
	if err := snapshots.Generate("http://127.0.0.1:8080/welcome", directory); err != nil {
		log.Fatal(err)
	}
	if err := snapshots.Generate("http://127.0.0.1:8080/todos", directory); err != nil {
		log.Fatal(err)
	}
	if err := files.CopyDirectory(filepath.Join("app", "dist", "assets"), filepath.Join(directory, "assets")); err != nil {
		log.Fatal(err)
	}
}
