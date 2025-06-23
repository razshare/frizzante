package frizzante

import (
	"github.com/razshare/frizzante/libglobals"
	"testing"
)

func TestKB(test *testing.T) {
	if 1024 != libglobals.KB {
		test.Fatalf("KB constant is not %d", 1024)
	}
}

func TestMB(test *testing.T) {
	expected := 1024 * 1024
	if expected != libglobals.MB {
		test.Fatalf("MB constant is not %d", expected)
	}
}

func TestGB(test *testing.T) {
	expected := 1024 * 1024 * 1024
	if expected != libglobals.GB {
		test.Fatalf("GB constant is not %d", expected)
	}
}

func TestTB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024
	if expected != libglobals.TB {
		test.Fatalf("TB constant is not %d", expected)
	}
}

func TestPB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024
	if expected != libglobals.PB {
		test.Fatalf("PB constant is not %d", expected)
	}
}

func TestEB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	if expected != libglobals.EB {
		test.Fatalf("EB constant is not %d", expected)
	}
}
