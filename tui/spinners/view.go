package spinners

import (
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func (model *Model) View() string {
	lines := strings.Split(model.Message, "\n")
	var builder strings.Builder
	for index, line := range lines {
		builder.WriteString("\r")
		if index == 0 {
			builder.WriteString(configs.Styles.Menu.Render(model.Spinner.View()))
		} else {
			builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("│"))
		}
		builder.WriteString(configs.Styles.Menu.Render(line))
		builder.WriteString("\n")
	}
	return builder.String()
	//return fmt.Sprintf("\r%s %s (lines=%d)\n", model.Spinner.View(), configs.Styles.Menu.Render(model.Message), len(lines))
}
