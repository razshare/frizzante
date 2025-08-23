package send

import "testing"

func TestStatus(t *testing.T) {
	client := MockClient()
	Status(client, 400)
	if client.Status != 400 {
		t.Fatal("status should be 400")
	}
}
