package messages

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/razshare/frizzante/tui/config"
)

func Status(label string, text string, labelBg string, labelFg string, color string) {
	labelWidth := 9

	labelStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(labelBg)).
		Foreground(lipgloss.Color(labelFg)).
		Bold(true).
		Width(labelWidth).
		Align(lipgloss.Center)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if i == 0 {
			fmt.Println(config.Styles.Menu.PaddingRight(1).Render("│") + labelStyle.Render(label) + " " + textStyle.Render(line))
		} else {
			fmt.Println(config.Styles.Menu.PaddingRight(1).Render("│") + labelStyle.Render("") + " " + textStyle.Render(line))
		}
	}
}
