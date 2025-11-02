package render

import (
	"embed"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

//go:generate rm -fr ./app
//go:generate mkdir -p ./app
//go:generate cp -r ../../../../app/dist ./app
//go:embed app
var TestNewEfs embed.FS

func TestNew(t *testing.T) {
	render := New()
	html, err := render(Options{Efs: TestNewEfs, View: views.View{Name: "Welcome"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Show Todos") {
		t.Fatal("view should contain Show Todos")
	}
}
