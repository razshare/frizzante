package security

import "testing"

func TestRandomHex(t *testing.T) {
	hex1 := RandomHex(16)

	if len(hex1) != 32 {
		t.Fatal("random hex should be 32 characters for 16 bytes")
	}

	hex2 := RandomHex(16)

	if hex1 == hex2 {
		t.Fatal("two random hex values should not be equal")
	}
}
