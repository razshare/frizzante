package send

import (
	"fmt"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func TestViewWithLocation(t *testing.T) {
	client := mock.NewClient()
	Header(client, "Location", "/about")
	View(client, views.View{}) // This should be a noop.
}

func TestViewWithAcceptJson(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Accept", "application/json")

	View(client, views.View{Name: "test", Props: map[string]any{"key": "value"}})

	writer := client.Writer.(*mock.ResponseWriter)

	if writer.MockHeader.Get("Cache-Control") != "no-store, no-cache, must-revalidate, max-age=0" {
		t.Fatal("cache control should be disabled")
	}

	if writer.MockHeader.Get("Pragma") != "no-cache" {
		t.Fatal("pragma should be no-cache")
	}

	if writer.MockHeader.Get("Content-Type") != "application/json" {
		t.Fatal("content type should be json")
	}

	if string(writer.MockBytes) != `{"name":"test","render":0,"align":0,"props":{"key":"value"}}` {
		t.Fatal("content should be view as json")
	}
}

func TestView(t *testing.T) {
	client := mock.NewClient()

	client.Config.Render = func(view views.View) (html string, err error) {

		return fmt.Sprintf("hello from %s", view.Name), nil
	}

	View(client, views.View{Name: "test", Props: map[string]any{"key": "value"}})

	writer := client.Writer.(*mock.ResponseWriter)

	if string(writer.MockBytes) != "hello from test" {
		t.Fatal("content should be hello from test")
	}
}
