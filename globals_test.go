package main

import (
	"github.com/razshare/frizzante/globals"
	"testing"
)

func TestKB(test *testing.T) {
	if 1024 != globals.KB {
		test.Fatalf("KB constant is not %d", 1024)
	}
}

func TestMB(test *testing.T) {
	expected := 1024 * 1024
	if expected != globals.MB {
		test.Fatalf("MB constant is not %d", expected)
	}
}

func TestGB(test *testing.T) {
	expected := 1024 * 1024 * 1024
	if expected != globals.GB {
		test.Fatalf("GB constant is not %d", expected)
	}
}

func TestTB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024
	if expected != globals.TB {
		test.Fatalf("TB constant is not %d", expected)
	}
}

func TestPB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024
	if expected != globals.PB {
		test.Fatalf("PB constant is not %d", expected)
	}
}

func TestEB(test *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	if expected != globals.EB {
		test.Fatalf("EB constant is not %d", expected)
	}
}
