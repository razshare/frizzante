package list

import "github.com/charmbracelet/lipgloss"

type Options struct {
	HeaderStyle    lipgloss.Style
	RowStyle       lipgloss.Style
	AltRowStyle    lipgloss.Style
	MaxColumnWidth int
	ColumnPadding  int
	HeaderHeight   int
}
