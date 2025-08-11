package messages

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func Status(label string, text string, bgColor string, fgColor string, textColor string) {
	labelWidth := 9

	labelStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor)).
		Bold(true).
		Width(labelWidth).
		Align(lipgloss.Center)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(textColor)).
		Bold(true)

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i == 0 {
			fmt.Println(labelStyle.Render(label) + " " + textStyle.Render(line))
		} else {
			fmt.Println(labelStyle.Render("") + " " + textStyle.Render(line))
		}
	}
}
