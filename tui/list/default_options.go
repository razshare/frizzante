package list

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/razshare/frizzante/tui/configs"
)

func DefaultOptions() Options {
	return Options{
		MaxColumnWidth: 60,
		ColumnPadding:  2,
		HeaderHeight:   2,
		HeaderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("221")),
		RowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(configs.Colors.Primary)),
		AltRowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(configs.Colors.Secondary)),
	}
}
