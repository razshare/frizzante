package spinners

import "github.com/charmbracelet/bubbles/spinner"

type Model struct {
	Spinner     spinner.Model
	Message     string
	QuitMessage string
}
