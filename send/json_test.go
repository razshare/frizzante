package send

import "testing"

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	client := MockClient()
	Json(client, Payload{Key: "value"})
	writer := client.Writer.(*MockWriter)
	if string(writer.MockBytes) != `{"key":"value"}` {
		t.Fatal("content should be json")
	}
}
