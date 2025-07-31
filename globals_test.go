package main

import (
	"github.com/razshare/frizzante/globals"
	"testing"
)

func TestKB(t *testing.T) {
	if 1024 != globals.KB {
		t.Fatalf("KB constant is not %d", 1024)
	}
}

func TestMB(t *testing.T) {
	expected := 1024 * 1024
	if expected != globals.MB {
		t.Fatalf("MB constant is not %d", expected)
	}
}

func TestGB(t *testing.T) {
	expected := 1024 * 1024 * 1024
	if expected != globals.GB {
		t.Fatalf("GB constant is not %d", expected)
	}
}

func TestTB(t *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024
	if expected != globals.TB {
		t.Fatalf("TB constant is not %d", expected)
	}
}

func TestPB(t *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024
	if expected != globals.PB {
		t.Fatalf("PB constant is not %d", expected)
	}
}

func TestEB(t *testing.T) {
	expected := 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	if expected != globals.EB {
		t.Fatalf("EB constant is not %d", expected)
	}
}
