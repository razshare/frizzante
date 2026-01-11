package send

import (
	"embed"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/render"
)

func TestViewWithLocation(t *testing.T) {
	client := mocks.NewClient()
	Header(client, "Location", "/about")
	View(client, views.View{}) // This should be a noop.
}

func TestViewWithAcceptJson(t *testing.T) {
	client := mocks.NewClient()
	client.Request.Header.Set("Accept", "application/json")

	View(client, views.View{Name: "test", Props: map[string]any{"key": "value"}})

	writer := client.Writer.(*mocks.ResponseWriter)

	if writer.MockHeader.Get("Cache-Control") != "no-store, no-cache, must-revalidate, max-age=0" {
		t.Fatal("cache control should be disabled")
	}

	if writer.MockHeader.Get("Pragma") != "no-cache" {
		t.Fatal("pragma should be no-cache")
	}

	if writer.MockHeader.Get("Content-Type") != "application/json" {
		t.Fatal("content type should be json")
	}

	if string(writer.MockBytes) != `{"isSnapshot":false,"name":"test","render":0,"align":0,"props":{"key":"value"}}` {
		t.Fatal("content should be view as json")
	}
}

//go:embed app/dist
var TestViewEfs embed.FS

func TestView(t *testing.T) {
	client := mocks.NewClient()
	client.Options.Efs = TestViewEfs

	var err error
	if client.Options.Render, err = render.New(); err != nil {
		t.Fatal(err)
	}

	View(client, views.View{Name: "Welcome", Props: map[string]any{"key": "value"}})

	writer := client.Writer.(*mocks.ResponseWriter)

	if !strings.Contains(string(writer.MockBytes), "Modern Go + Svelte Framework") {
		t.Fatal("view should contain Modern Go + Svelte Framework")
	}
}
