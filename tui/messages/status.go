package messages

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func Status(lbl string, txt string, bgc string, fgc string, txtc string) {
	labelWidth := 9

	labelStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(bgc)).
		Foreground(lipgloss.Color(fgc)).
		Bold(true).
		Width(labelWidth).
		Align(lipgloss.Center)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(txtc)).
		Bold(true)

	lines := strings.Split(txt, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if i == 0 {
			fmt.Println(labelStyle.Render(lbl) + " " + textStyle.Render(line))
		} else {
			fmt.Println(labelStyle.Render("") + " " + textStyle.Render(line))
		}
	}
}
