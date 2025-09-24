package wrap

import (
	"reflect"
	"testing"
)

func TestSend(t *testing.T) {
	type TestData struct {
		name     string
		text     string
		width    int
		expected []string
	}

	data := []TestData{
		{
			name:     "empty text",
			text:     "",
			width:    10,
			expected: nil,
		},
		{
			name:     "single line within width",
			text:     "hello",
			width:    10,
			expected: []string{"hello"},
		},
		{
			name:     "single line exceeds width",
			text:     "hello world test",
			width:    10,
			expected: []string{"hello", "world test"},
		},
		{
			name:     "multiple words wrap",
			text:     "this is a long line that needs wrapping",
			width:    15,
			expected: []string{"this is a long", "line that needs", "wrapping"},
		},
		{
			name:     "preserve newlines",
			text:     "line one\nline two\nline three",
			width:    20,
			expected: []string{"line one", "line two", "line three"},
		},
		{
			name:     "empty lines skipped",
			text:     "line one\n\nline three",
			width:    20,
			expected: []string{"line one", "line three"},
		},
		{
			name:     "zero width returns split lines",
			text:     "line one\nline two",
			width:    0,
			expected: []string{"line one", "line two"},
		},
		{
			name:     "negative width returns split lines",
			text:     "line one\nline two",
			width:    -1,
			expected: []string{"line one", "line two"},
		},
		{
			name:     "single long word",
			text:     "superlongwordthatcannotbewrapped",
			width:    10,
			expected: []string{"superlongwordthatcannotbewrapped"},
		},
		{
			name:     "mixed content with newlines",
			text:     "short\nthis is a longer line that needs wrapping\nlast",
			width:    20,
			expected: []string{"short", "this is a longer", "line that needs", "wrapping", "last"},
		},
		{
			name:     "exact width boundary",
			text:     "hello world",
			width:    11,
			expected: []string{"hello world"},
		},
		{
			name:     "exact width boundary plus one",
			text:     "hello worlds",
			width:    11,
			expected: []string{"hello", "worlds"},
		},
		{
			name:     "multiple spaces between words",
			text:     "hello     world",
			width:    20,
			expected: []string{"hello     world"}, // preserves spaces when within width
		},
		{
			name:     "leading and trailing spaces",
			text:     "  hello world  ",
			width:    20,
			expected: []string{"  hello world  "}, // preserves spaces when within width
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := Send(d.text, d.width)
			if !reflect.DeepEqual(result, d.expected) {
				t.Errorf("Send() = %v, want %v", result, d.expected)
			}
		})
	}
}
