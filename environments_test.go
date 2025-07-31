package main

import (
	"github.com/razshare/frizzante/environments"
	"os"
	"testing"
)

func TestLoadDotenv(test *testing.T) {
	err := environments.LoadDotenv("test.env")
	if err != nil {
		test.Fatal(err)
	}

	name := os.Getenv("name")

	if name != "world" {
		test.Fatal("environment variable name is not \"world\"")
	}
}
