package security

import (
	"strings"
	"testing"
)

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

func TestRandomBase64(t *testing.T) {
	b64 := RandomBase64(24)

	if len(b64) == 0 {
		t.Fatal("random base64 should not be empty")
	}

	if strings.Contains(b64, "+") || strings.Contains(b64, "/") {
		t.Fatal("URL-safe base64 should not contain + or /")
	}
}

func TestRandomBase64Standard(t *testing.T) {
	b64 := RandomBase64Standard(24)

	if len(b64) == 0 {
		t.Fatal("random base64 should not be empty")
	}
}
