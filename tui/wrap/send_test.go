package wrap

import (
	"slices"
	"testing"
)

func TestSend(t *testing.T) {
	var lines []string

	// empty text
	lines = Send("", 5)
	if !slices.Equal(lines, []string{}) {
		t.Fatal("wrap should be empty")
	}

	// exceeding
	lines = Send("some text", 5)
	if !slices.Equal(lines, []string{"some", "text"}) {
		t.Fatal("wrap should contain some and text")
	}

	// preserve newlines
	lines = Send("line one\nline two\nline three", 20)
	if !slices.Equal(lines, []string{"line one", "line two", "line three"}) {
		t.Fatal("wrap should contain line one, line two and line three")
	}

	// empty lines skipped
	lines = Send("line one\n\nline three", 20)
	if !slices.Equal(lines, []string{"line one", "line three"}) {
		t.Fatal("wrap should contain line one and line three")
	}

	// zero width returns split lines
	lines = Send("line one\nline two", 0)
	if !slices.Equal(lines, []string{"line one", "line two"}) {
		t.Fatal("wrap should contain line one and line two")
	}

	// single long word
	lines = Send("superlongwordthatcannotbewrapped", 10)
	if !slices.Equal(lines, []string{"superlongwordthatcannotbewrapped"}) {
		t.Fatal("wrap should contain superlongwordthatcannotbewrapped")
	}

	// mixed content with newlines
	lines = Send("short\nthis is a longer line that needs wrapping\nlast", 20)
	if !slices.Equal(lines, []string{"short", "this is a longer", "line that needs", "wrapping", "last"}) {
		t.Fatal("wrap should contain 5 lines")
	}

	// exact width boundary
	lines = Send("hello world", 11)
	if !slices.Equal(lines, []string{"hello world"}) {
		t.Fatal("wrap should contain hello world")
	}

	// exact width boundary plus one
	lines = Send("hello worlds", 11)
	if !slices.Equal(lines, []string{"hello", "worlds"}) {
		t.Fatal("wrap should contain hello and worlds")
	}

	// multiple spaces between words
	lines = Send("hello     world", 20)
	if !slices.Equal(lines, []string{"hello     world"}) {
		t.Fatal("wrap should contain hello     world")
	}

	// leading and trailing spaces
	lines = Send("  hello world  ", 20)
	if !slices.Equal(lines, []string{"  hello world  "}) {
		t.Fatal("wrap should contain   hello world  ")
	}
}
