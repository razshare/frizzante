package config

import "github.com/charmbracelet/lipgloss"

type ThemeColors struct {
	Primary   string
	Secondary string
	Success   string
	Error     string
	Warning   string
	Info      string
	Tip       string
	Muted     string
}

type ThemeStyles struct {
	Title     lipgloss.Style
	Popup     lipgloss.Style
	Item      lipgloss.Style
	Selected  lipgloss.Style
	Status    func(color string) lipgloss.Style
	BigText   lipgloss.Style
	Section   lipgloss.Style
	Subheader lipgloss.Style
	Menu      lipgloss.Style
	UserGuide lipgloss.Style
	UserInput lipgloss.Style
	Spinner   lipgloss.Style
	Flag      lipgloss.Style
	Category  lipgloss.Style
	Example   lipgloss.Style
}
