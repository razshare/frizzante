package ssr

import (
	"embed"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

//go:generate rm -fr ./app
//go:generate mkdir -p ./app
//go:generate cp -r ../../../../app/dist ./app
//go:embed app
var TestNewEfs embed.FS

func TestNew(t *testing.T) {
	f := New(Config{Efs: TestNewEfs})
	html, err := f(view.View{Name: "Welcome"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Powered by Svelte for smooth interfaces") {
		t.Fatal("view should contain Powered by Svelte for smooth interfaces")
	}
}
