package security

import "testing"

func TestSha1(t *testing.T) {
	hash := Sha1("hello")

	if hash != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Fatal("hash should be aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d")
	}
}
