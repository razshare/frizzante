package cli

import (
	"embed"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FeatureCopyInstruction struct {
	FeatureName          string
	OriginDirectory      string
	DestinationDirectory string
}

type Cli struct {
	Efs    embed.FS
	Parsed bool
}

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
	Title       lipgloss.Style
	Item        lipgloss.Style
	Selected    lipgloss.Style
	Status      func(color string) lipgloss.Style
	BigText     lipgloss.Style
	Section     lipgloss.Style
	Subheader   lipgloss.Style
	Spinner     lipgloss.Style
	Flag        lipgloss.Style
	Category    lipgloss.Style
	Example     lipgloss.Style
}

type ConfirmModel struct {
	Prompt       string
	Confirmed    bool
	DefaultValue bool
}

type InputModel struct {
	TextInput textinput.Model
	Prompt    string
	Cancelled bool
}

type MultiSelectModel struct {
	Choices  []string
	Cursor   int
	Selected map[int]bool
	Prompt   string
}

// SimpleChooseModel defines single selection options
type SimpleChooseModel struct {
	Choices         []string        // All available Choices
	FilteredChoices []string        // Choices after filtering
	Cursor          int             // Current Cursor position in filtered list
	Prompt          string          // The Prompt to display
	Selected        string          // The Selected choice
	SearchInput     textinput.Model // Search input field
	Searching       bool            // Whether we're in search mode
	ViewportStart   int             // Start index for viewport (for scrolling)
	MaxVisible      int             // Maximum visible items (5 by default)
}

type SpinnerModel struct {
	Spinner spinner.Model
	Message string
}

type SpinnerManager struct {
	Model   SpinnerModel
	Program *tea.Program
	Done    chan bool
}
