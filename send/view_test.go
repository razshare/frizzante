package send

import (
	"fmt"
	"github.com/razshare/frizzante/mock"
	"github.com/razshare/frizzante/view"
	"testing"
)

func TestViewWithLocation(t *testing.T) {
	c := mock.NewClient()
	Header(c, "Location", "/about")
	View(c, view.View{}) // This should be a noop.
}

func TestViewWithAcceptJson(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("Accept", "application/json")

	View(c, view.View{Name: "test", Props: map[string]any{"key": "value"}})

	w := c.Writer.(*mock.ResponseWriter)

	if w.MockHeader.Get("Cache-Control") != "no-store, no-cache, must-revalidate, max-age=0" {
		t.Fatal("cache control should be disabled")
	}

	if w.MockHeader.Get("Pragma") != "no-cache" {
		t.Fatal("pragma should be no-cache")
	}

	if w.MockHeader.Get("Content-Type") != "application/json" {
		t.Fatal("content type should be json")
	}

	if string(w.MockBytes) != `{"align":0,"name":"test","props":{"key":"value"},"render":0}` {
		t.Fatal("content should be view as json")
	}
}

func TestView(t *testing.T) {
	c := mock.NewClient()
	c.Config.Render = func(v view.View) (html string, err error) {
		return fmt.Sprintf("hello from %s", v.Name), nil
	}

	View(c, view.View{Name: "test", Props: map[string]any{"key": "value"}})

	w := c.Writer.(*mock.ResponseWriter)

	if string(w.MockBytes) != "hello from test" {
		t.Fatal("content should be hello from test")
	}
}
