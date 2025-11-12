package spinners

import (
	"fmt"

	"github.com/razshare/frizzante/tui/configs"
)

func (model *Model) View() string {
	return fmt.Sprintf("\r%s %s", model.Spinner.View(), configs.Styles.Menu.Render(model.Message))
}
