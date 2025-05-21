package pages

import "testing"

func TestController(test *testing.T) {
	pageName.WithViewRoot("templates")
	actual := pageName.FindId()

	expected := "pages"
	if actual != expected {
		test.Fatalf("package id was expected to be `%s`, but received `%s` instead", expected, actual)
	}
}
