package spinners

import (
	"fmt"

	"github.com/razshare/frizzante/tui/config"
)

func (model *Model) View() string {
	return fmt.Sprintf("\r%s %s", model.Spinner.View(), config.Styles.Menu.Render(model.Message))
}
