package send

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestStatus(t *testing.T) {
	c := mock.NewClient()
	Status(c, 400)
	if c.Status != 400 {
		t.Fatal("status should be 400")
	}
}
