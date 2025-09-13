package ssr

import (
	"embed"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

//go:embed app
var EfsTestNew embed.FS

func TestNew(t *testing.T) {
	f := New(Config{Efs: EfsTestNew})
	html, err := f(view.View{Name: "Welcome"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Powered by Svelte for smooth interfaces") {
		t.Fatal("view should contain Powered by Svelte for smooth interfaces")
	}
}
