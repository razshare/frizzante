package table

import "github.com/charmbracelet/lipgloss"

type Options struct {
	MaxColumnWidth int
	ColumnPadding  int
	HeaderHeight   int
	HeaderStyle    lipgloss.Style
	RowStyle       lipgloss.Style
	AltRowStyle    lipgloss.Style
}
