package frizzante

import (
	"testing"
)

func TestKB(test *testing.T) {
	if 1024 != KB {
		test.Fatalf("KB constant is not %d", 1024)
	}
}

func TestMB(test *testing.T) {
	expected := 1024 * 1024
	if expected != MB {
		test.Fatalf("MB constant is not %d", expected)
	}
}

func TestGB(test *testing.T) {
	expected := 1024 * 1024 * 1024
	if expected != GB {
		test.Fatalf("GB constant is not %d", expected)
	}
}

func TestTB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024
	if expected != TB {
		test.Fatalf("TB constant is not %d", expected)
	}
}

func TestPB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024
	if expected != PB {
		test.Fatalf("PB constant is not %d", expected)
	}
}

func TestEB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	if expected != EB {
		test.Fatalf("EB constant is not %d", expected)
	}
}
