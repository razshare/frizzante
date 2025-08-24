package ssr

import (
	"embed"
	"github.com/razshare/frizzante/view"
	"strings"
	"testing"
)

//go:embed app
var EfsTestNew embed.FS

func TestNew(t *testing.T) {
	f := New(Config{Efs: EfsTestNew})
	html, err := f(view.View{Name: "Welcome"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<h1>Welcome to Frizzante.</h1>") {
		t.Fatal("view should container <h1>Welcome to Frizzante.</h1>")
	}
}
