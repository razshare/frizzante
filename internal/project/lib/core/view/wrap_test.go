package view

import "testing"

func TestData(t *testing.T) {
	data := Wrap(View{
		Name:  "name",
		Title: "title",
		Props: map[string]any{
			"key": "value",
		},
	})

	if data.Props.(map[string]any)["key"] != "value" {
		t.Fatal("key should be value")
	}
}
