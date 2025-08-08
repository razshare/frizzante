package config

import "github.com/charmbracelet/lipgloss"

type ThemeColors struct {
	Primary   string
	Secondary string
	Success   string
	Error     string
	Warning   string
	Info      string
	Muted     string
}

type ThemeStyles struct {
	Title     lipgloss.Style
	Item      lipgloss.Style
	Selected  lipgloss.Style
	Status    func(color string) lipgloss.Style
	BigText   lipgloss.Style
	Section   lipgloss.Style
	Subheader lipgloss.Style
	Spinner   lipgloss.Style
	Flag      lipgloss.Style
	Category  lipgloss.Style
	Example   lipgloss.Style
}
