package types

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func Generate[T any]() {
	var val T
	var err error
	var primary string
	var secondary string

	t := reflect.TypeOf(val)

	if primary, secondary, _, err = Extract(t, make([]string, 0)); err != nil {
		log.Fatal(err)
	}

	if !files.IsDirectory(filepath.Join(".gen", "types")) {
		if err = os.MkdirAll(filepath.Join(".gen", "types"), os.ModePerm); err != nil {
			return
		}
	}

	befores := []string{
		"github.com/razshare/frizzante/internal/project",
		"github.com/razshare/frizzante/internal/additions",
	}
	after := "main"
	pkg := t.PkgPath()

	for _, before := range befores {
		pkg = strings.ReplaceAll(pkg, before, after)
	}

	dname := filepath.Join(".gen", "types", strings.ReplaceAll(pkg, "/", string(filepath.Separator)))
	if !files.IsDirectory(dname) {
		if err = os.MkdirAll(dname, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	fname := filepath.Join(dname, t.Name()+".d.ts")
	if err = os.WriteFile(fname, []byte(primary+secondary), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
