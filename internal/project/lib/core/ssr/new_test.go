package ssr

import (
	"embed"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

//go:generate rm -fr ./app
//go:generate mkdir -p ./app
//go:generate cp -r ../../../../app/dist ./app
//go:embed app
var TestNewEfs embed.FS

func TestNewFunction(t *testing.T) {
	var html string
	var err error
	render := New(1)
	if html, err = render(renders.RenderOptions{
		Efs: TestNewEfs,
		View: views.View{
			Name: "Welcome",
		},
	}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Show Todos") {
		t.Fatal("view should contain Show Todos")
	}
}
