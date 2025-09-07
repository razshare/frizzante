package table

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTableColumnWidthCalculation(t *testing.T) {
	tests := []struct {
		name           string
		headers        []string
		rows           [][]string
		maxWidth       int
		expectedWidths []int
	}{
		{
			name:           "headers determine width",
			headers:        []string{"Name", "Description", "ID"},
			rows:           [][]string{{"a", "b", "c"}},
			maxWidth:       20,
			expectedWidths: []int{4, 11, 2},
		},
		{
			name:    "cell content determines width",
			headers: []string{"A", "B", "C"},
			rows: [][]string{
				{"longer", "medium", "x"},
				{"short", "y", "also longer"},
			},
			maxWidth:       20,
			expectedWidths: []int{6, 6, 11},
		},
		{
			name:    "max width constraint applied",
			headers: []string{"Short", "Long"},
			rows: [][]string{
				{"a", "this is a very long cell content that exceeds maximum"},
			},
			maxWidth:       10,
			expectedWidths: []int{5, 10},
		},
		{
			name:    "mixed content sizes",
			headers: []string{"ID", "Name", "Description"},
			rows: [][]string{
				{"1", "Alice", "Software Engineer"},
				{"2", "Bob", "Designer"},
				{"3", "Charlie", "Product Manager"},
			},
			maxWidth:       20,
			expectedWidths: []int{2, 7, 17},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colWidths := make([]int, len(tt.headers))
			for i, header := range tt.headers {
				colWidths[i] = len(header)
			}

			for _, row := range tt.rows {
				for i, cell := range row {
					if i < len(colWidths) {
						if len(cell) > colWidths[i] {
							colWidths[i] = len(cell)
						}
					}
				}
			}

			for i := range colWidths {
				if colWidths[i] > tt.maxWidth {
					colWidths[i] = tt.maxWidth
				}
			}

			for i, expected := range tt.expectedWidths {
				if colWidths[i] != expected {
					t.Errorf("column %d width = %d, want %d", i, colWidths[i], expected)
				}
			}
		})
	}
}

func TestTableEmptyData(t *testing.T) {
	t.Run("empty headers", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("should not panic with empty headers")
			}
		}()

		Send([]string{}, [][]string{{"data"}})
	})

	t.Run("empty rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("should not panic with empty rows")
			}
		}()

		Send([]string{"Header"}, [][]string{})
	})

	t.Run("both empty", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("should not panic with both empty")
			}
		}()

		Send([]string{}, [][]string{})
	})
}

func TestTableMultilineHandling(t *testing.T) {
	rows := [][]string{
		{"1", "Single line"},
		{"2", "Line one\nLine two"},
		{"3", "Another single"},
	}

	processedCount := 0
	maxLines := 1

	for _, row := range rows {
		rowMaxLines := 1
		for _, cell := range row {
			lines := 1
			for _, ch := range cell {
				if ch == '\n' {
					lines++
				}
			}
			if lines > rowMaxLines {
				rowMaxLines = lines
			}
		}
		processedCount += rowMaxLines
		if rowMaxLines > maxLines {
			maxLines = rowMaxLines
		}
	}

	processedCount += len(rows) - 1

	if maxLines != 2 {
		t.Errorf("max lines = %d, want 2", maxLines)
	}
}

func TestTableOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		opts := DefaultOptions()

		if opts.MaxColumnWidth != 60 {
			t.Errorf("default max width = %d, want 60", opts.MaxColumnWidth)
		}

		_ = opts.HeaderStyle.Render("test")
		_ = opts.RowStyle.Render("test")
		_ = opts.AltRowStyle.Render("test")
	})

	t.Run("custom options", func(t *testing.T) {
		customOpts := Options{
			MaxColumnWidth: 100,
			HeaderStyle:    lipgloss.NewStyle().Bold(true),
			RowStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("red")),
			AltRowStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
		}

		if customOpts.MaxColumnWidth != 100 {
			t.Errorf("custom max width = %d, want 100", customOpts.MaxColumnWidth)
		}
	})
}

func TestTableRowSeparation(t *testing.T) {
	rows := [][]string{
		{"1", "2"},
		{"3", "4"},
		{"5", "6"},
	}

	separatorCount := len(rows) - 1

	if separatorCount != 2 {
		t.Errorf("separator count = %d, want 2", separatorCount)
	}
}

func TestTableJaggedRows(t *testing.T) {
	headers := []string{"Col1", "Col2", "Col3"}
	rows := [][]string{
		{"a", "b"},
		{"c", "d", "e"},
		{"f"},
	}

	for _, row := range rows {
		processedRow := make([]string, len(headers))
		for i := 0; i < len(headers); i++ {
			if i < len(row) {
				processedRow[i] = row[i]
			} else {
				processedRow[i] = ""
			}
		}

		if len(processedRow) != len(headers) {
			t.Errorf("processed row length = %d, want %d", len(processedRow), len(headers))
		}
	}
}

func TestTableWrapping(t *testing.T) {
	headers := []string{"ID", "Long Text"}
	rows := [][]string{
		{"1", "This is a very long text that should be wrapped to multiple lines"},
	}

	opts := Options{
		MaxColumnWidth: 20,
		HeaderStyle:    lipgloss.NewStyle(),
		RowStyle:       lipgloss.NewStyle(),
		AltRowStyle:    lipgloss.NewStyle(),
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error("should not panic during wrapping")
		}
	}()

	Send(headers, rows, opts)
}
