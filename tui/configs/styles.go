package configs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/razshare/frizzante/tui/themes"
)

var Styles = themes.Styles{
	Title: lipgloss.NewStyle().Bold(true),

	Item: lipgloss.NewStyle(),

	Popup: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Info)),

	Selected: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)),

	Status: func(color string) lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)
	},

	BigText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)).
		Bold(true).
		Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(Colors.Primary)),

	Section: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)).
		Bold(true).
		Underline(true),

	Subheader: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Warning)).
		Bold(true),

	Menu: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)),

	UserInput: lipgloss.NewStyle(),

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
		Foreground(lipgloss.Color(Colors.Input)),

	InfoLabel: lipgloss.
		NewStyle().
		Background(lipgloss.Color(Colors.Info)).
		Foreground(lipgloss.Color("17")),

	InfoText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Info)),

	WarningLabel: lipgloss.
		NewStyle().
		Background(lipgloss.Color(Colors.Warning)).
		Foreground(lipgloss.Color("17")),

	WarningText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Warning)),

	ErrorLabel: lipgloss.
		NewStyle().
		Background(lipgloss.Color(Colors.Error)).
		Foreground(lipgloss.Color("17")),

	ErrorText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Error)),
}
