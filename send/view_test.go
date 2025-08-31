package send

import (
	"fmt"
	"testing"

	"github.com/razshare/frizzante/mock"
	"github.com/razshare/frizzante/view"
)

func TestViewWithLocation(t *testing.T) {
	client := mock.NewClient()
	Header(client, "Location", "/about")
	View(client, view.View{}) // This should be a noop.
}

func TestViewWithAcceptJson(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Accept", "application/json")

	View(client, view.View{Name: "test", Props: map[string]any{"key": "value"}})

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

	if string(writer.MockBytes) != `{"align":0,"name":"test","props":{"key":"value"},"render":0}` {
		t.Fatal("content should be view as json")
	}
}

func TestView(t *testing.T) {
	client := mock.NewClient()
	client.Config.Render = func(v view.View) (html string, err error) {
		return fmt.Sprintf("hello from %s", v.Name), nil
	}

	View(client, view.View{Name: "test", Props: map[string]any{"key": "value"}})

	writer := client.Writer.(*mock.ResponseWriter)

	if string(writer.MockBytes) != "hello from test" {
		t.Fatal("content should be hello from test")
	}
}
