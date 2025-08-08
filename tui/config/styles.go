package config

import "github.com/charmbracelet/lipgloss"

var Styles = ThemeStyles{
	Title: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)).
		Bold(true),

	Item: lipgloss.NewStyle().
		PaddingLeft(2),

	Selected: lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(lipgloss.Color(Colors.Secondary)),

	Status: func(color string) lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			PaddingLeft(2)
	},

	BigText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Info)).
		Bold(true).
		Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(Colors.Primary)).
		Padding(1, 4).
		Margin(1, 2),

	Section: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)).
		Bold(true).
		Underline(true).
		Padding(1, 0),

	Subheader: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Warning)).
		Bold(true).
		Padding(1, 0),

	Spinner: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Info)),

	Flag: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)).Bold(true),

	Category: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)).Bold(true).Underline(true),

	Example: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Muted)),
}
