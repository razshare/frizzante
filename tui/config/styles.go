package config

import "github.com/charmbracelet/lipgloss"

var Styles = ThemeStyles{
	Title: lipgloss.NewStyle().Bold(true),

	Item: lipgloss.NewStyle().
		PaddingLeft(1),

	Popup: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Info)),

	Selected: lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(lipgloss.Color(Colors.Primary)),

	Status: func(color string) lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			PaddingLeft(1)
	},

	BigText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)).
		Bold(true).
		Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(Colors.Primary)).
		Padding(1, 1),

	Section: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)).
		Bold(true).
		Underline(true).
		Padding(1, 0),

	Subheader: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Warning)).
		Bold(true).
		Padding(1, 0),

	Menu: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)),

	UserInput: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Muted)),

	UserGuide: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Muted)),

	Spinner: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)),

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
