package main

import (
	"github.com/razshare/frizzante/environments"
	"os"
	"testing"
)

func TestLoadDotenv(t *testing.T) {
	err := environments.LoadDotenv()
	if err != nil {
		t.Fatal(err)
	}

	name := os.Getenv("name")

	if name != "world" {
		t.Fatal("environment variable name is not \"world\"")
	}
}
