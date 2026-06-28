package views

import "testing"

func TestData(t *testing.T) {
	data := NewData(View{
		Name: "name",
		Props: map[string]any{
			"key": "value",
		},
	})
	if data.Props.(map[string]any)["key"] != "value" {
		t.Fatal("key should be value")
	}
}
