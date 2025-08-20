package npmselect

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/program"
)

func Send() ([]string, error) {
	model, err := program.Run(&Model{
		Search:        initSearch(),
		Viewport:      initViewport(),
		Packages:      []PackageInfo{},
		Selected:      []string{},
		Loading:       false,
		DebounceTimer: nil,
	})

	if err != nil {
		if errors.Is(err, tea.ErrInterrupted) {
			return nil, err
		}
		return nil, err
	}

	if model.Quitting && !model.Confirmed {
		return nil, errors.New("cancelled")
	}

	return ExtractPackageNames(model.Selected), nil
}