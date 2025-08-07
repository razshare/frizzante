package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

var Colors = ThemeColors{
	Primary:   "87",
	Secondary: "169",
	Success:   "82",
	Error:     "196",
	Warning:   "214",
	Info:      "39",
	Muted:     "240",
}

var Styles = ThemeStyles{
	Title: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)).
		Bold(true),

	Item: lipgloss.NewStyle().
		PaddingLeft(2),

	Selected: lipgloss.NewStyle().
		PaddingLeft(1).
		Foreground(lipgloss.Color(Colors.Secondary)),

	Status: func(color string) lipgloss.Style {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			PaddingLeft(2)
	},

	BigText: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Info)).
		Bold(true).
		Align(lipgloss.Center).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(Colors.Primary)).
		Padding(1, 4).
		Margin(1, 2),

	Section: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)).
		Bold(true).
		Underline(true).
		Padding(1, 0),

	Subheader: lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.Warning)).
		Bold(true).
		Padding(1, 0),

	Spinner: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Info)),

	Flag: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Secondary)).Bold(true),

	Category: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Primary)).Bold(true).Underline(true),

	Example: lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(Colors.Muted)),
}

func RunProgram[T tea.Model](model T) (T, error) {
	result, err := tea.NewProgram(model, tea.WithFPS(120)).Run()
	if err != nil {
		return model, err
	}
	if typed, ok := result.(T); ok {
		return typed, nil
	}
	return model, fmt.Errorf("unexpected model type")
}

func CharmChoose(prompt string, options []string) (string, error) {
	// Initialize the search input
	searchInput := textinput.New()
	searchInput.Width = 80

	result, err := RunProgram(&SimpleChooseModel{
		Choices:         options,
		FilteredChoices: options,
		Prompt:          prompt,
		SearchInput:     searchInput,
		MaxVisible:      6,
		ViewportStart:   0,
		Cursor:          0,
		Searching:       false,
	})
	if err != nil {
		return "", err
	}
	if result.Selected == "" {
		return "", fmt.Errorf("no selection made")
	}
	return result.Selected, nil
}

func CharmMultiSelect(prompt string, options []string) ([]string, error) {
	m := MultiSelectModel{Choices: options, Selected: make(map[int]bool), Prompt: prompt}
	result, err := RunProgram(m)
	if err != nil {
		return nil, err
	}
	var selections []string
	for i, choice := range result.Choices {
		if result.Selected[i] {
			selections = append(selections, choice)
		}
	}
	return selections, nil
}

func CharmInput(prompt string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.Width = 50
	m := InputModel{TextInput: ti, Prompt: prompt}
	result, err := RunProgram(m)
	if err != nil {
		return "", err
	}
	if result.Cancelled {
		return "", fmt.Errorf("input cancelled")
	}
	return result.TextInput.Value(), nil
}

func CharmConfirm(prompt string, defaultValue bool) (bool, error) {
	model := ConfirmModel{Prompt: prompt, DefaultValue: defaultValue, Confirmed: defaultValue}
	result, err := RunProgram(model)
	if err != nil {
		return false, err
	}
	return result.Confirmed, nil
}

func PrintStatusMessage(label string, text string, bgColor string, fgColor string, textColor string) {
	labelWidth := 9
	
	if label == "INFO" {
		labelWidth = 8
	}
	
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

func CharmSuccess(text string) {
	PrintStatusMessage("SUCCESS", text, Colors.Success, "0", Colors.Success)
}

func CharmError(text string) {
	PrintStatusMessage("ERROR", text, Colors.Error, "15", Colors.Error)
}

func CharmWarning(text string) {
	PrintStatusMessage("WARNING", text, Colors.Warning, "0", Colors.Warning)
}

func CharmInfo(text string) {
	PrintStatusMessage("INFO", text, Colors.Info, "15", Colors.Info)
}

func CharmSection(text string) {
	fmt.Println(Styles.Section.Render("## " + text))
}

func CharmSubheader(text string) {
	fmt.Println(Styles.Subheader.Render(text))
}

func CharmDockerHelp() {
	fmt.Println(Styles.Title.Render("🐙 You're running Frizzante in Docker!"))

	fmt.Println(Styles.Subheader.Render("⚡ Simple workflow:"))

	fmt.Println(Styles.Status(Colors.Info).Render("• Attach to the container: ") +
		Styles.Example.Render("docker exec -it frizzante-start sh"))

	fmt.Println(Styles.Status(Colors.Info).Render("• Run environment in container:"))
	fmt.Println(Styles.Item.Render("    • Dev environment: ") + Styles.Flag.Render("make dev"))
	fmt.Println(Styles.Item.Render("    • Prod environment: ") + Styles.Flag.Render("make build"))
	fmt.Println(Styles.Item.Render("    • To run the app: ") + Styles.Flag.Render("./.gen/bin/app"))

	fmt.Println(Styles.Status(Colors.Info).Render("• Run prod via docker:"))
	fmt.Println(Styles.Item.Render("    • Build image: ") +
		Styles.Example.Render("docker build --target frizzante_prod -t my-app:prod ."))
	fmt.Println(Styles.Item.Render("    • Run image: ") +
		Styles.Example.Render("docker run -p 8080:8080 my-app:prod"))
	fmt.Println(Styles.Item.Render("    • Via docker compose: ") +
		Styles.Example.Render("docker compose -f compose.yaml -f compose.prod.yaml up -d --build"))

	fmt.Println(Styles.Subheader.Render("🎉 Enjoy!!"))
	fmt.Println()
}

func WrapText(text string, width int) []string {
	if width <= 0 || len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := ""
	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func CharmTable(headers []string, rows [][]string) {
	columns := make([]table.Column, len(headers))
	maxColWidth := 60

	for index, header := range headers {
		width := len(header)
		for _, row := range rows {
			if index < len(row) && len(row[index]) > width {
				width = len(row[index])
			}
		}
		if width > maxColWidth {
			width = maxColWidth
		}
		columns[index] = table.Column{Title: header, Width: width + 2}
	}

	wrappedRows := make([]table.Row, 0)
	for rowIdx, row := range rows {
		maxLines := 1
		wrappedCells := make([][]string, len(row))

		for i, cell := range row {
			if i < len(columns) {
				cellWidth := columns[i].Width - 2
				wrappedCells[i] = WrapText(cell, cellWidth)
				if len(wrappedCells[i]) > maxLines {
					maxLines = len(wrappedCells[i])
				}
			}
		}

		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			newRow := make([]string, len(row))
			for cellIdx, wrappedCell := range wrappedCells {
				if lineIdx < len(wrappedCell) {
					newRow[cellIdx] = wrappedCell[lineIdx]
				} else {
					newRow[cellIdx] = ""
				}
			}
			wrappedRows = append(wrappedRows, newRow)
		}

		// Add empty row after each logical row (except the last one)
		if rowIdx < len(rows)-1 {
			emptyRow := make([]string, len(row))
			for index := range emptyRow {
				emptyRow[index] = ""
			}
			wrappedRows = append(wrappedRows, emptyRow)
		}
	}

	tableLocal := table.New(
		table.WithColumns(columns),
		table.WithRows(wrappedRows),
		table.WithFocused(false),
		table.WithHeight(len(wrappedRows)+2),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	// Remove selection highlight
	styles.Selected = lipgloss.NewStyle()
	tableLocal.SetStyles(styles)

	fmt.Println(tableLocal.View())
}

func CharmSpinner(message string) *SpinnerManager {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = Styles.Spinner

	m := SpinnerModel{
		Spinner: s,
		Message: message,
	}

	return &SpinnerManager{
		Model:   m,
		Program: tea.NewProgram(m),
	}
}

func (manager *SpinnerManager) Start() (err error) {
	manager.Done = make(chan bool)
	go func() {
		_, err = manager.Program.Run()
		close(manager.Done)
	}()
	time.Sleep(100 * time.Millisecond)
	return
}

func (manager *SpinnerManager) Stop() {
	manager.Program.Quit()
	<-manager.Done
	fmt.Print("\r\033[K")
}

func (model SpinnerModel) Init() tea.Cmd {
	return model.Spinner.Tick
}

func (model SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model.Spinner, cmd = model.Spinner.Update(msg)
	return model, cmd
}

func (model SpinnerModel) View() string {
	return fmt.Sprintf("%s %s", model.Spinner.View(), model.Message)
}

func (model *SimpleChooseModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model *SimpleChooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		key := messageLocal.String()
		switch key {
		case "ctrl+c":
			return model, tea.Quit
		case "esc":
			if model.Searching {
				model.ResetSearch()
			}
			return model, nil
		case "enter":
			if len(model.FilteredChoices) > 0 {
				model.Selected = model.FilteredChoices[model.Cursor]
				return model, tea.Quit
			}
		case "up", "ctrl+p":
			model.Navigate(-1)
			return model, nil
		case "down", "ctrl+n", "tab":
			model.Navigate(1)
			return model, nil
		case "backspace", "ctrl+h":
			if model.Searching {
				return model.HandleSearchInput(messageLocal)
			}
		default:
			if len(key) == 1 || key == "space" {
				if !model.Searching {
					model.Searching = true
					model.SearchInput.Focus()
				}
				return model.HandleSearchInput(messageLocal)
			}
		}
	}
	return model, nil
}

func (model *SimpleChooseModel) Navigate(direction int) {
	count := len(model.FilteredChoices)
	if count == 0 {
		return
	}
	
	model.Cursor = (model.Cursor + direction + count) % count
	
	if model.Cursor < model.ViewportStart {
		model.ViewportStart = model.Cursor
	} else if model.Cursor >= model.ViewportStart+model.MaxVisible {
		model.ViewportStart = model.Cursor - model.MaxVisible + 1
	}
}

func (model *SimpleChooseModel) ResetSearch() {
	model.SearchInput.SetValue("")
	model.Searching = false
	model.FilteredChoices = model.Choices
	model.Cursor = 0
	model.ViewportStart = 0
}

func (model *SimpleChooseModel) HandleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prevValue := model.SearchInput.Value()
	var cmd tea.Cmd
	model.SearchInput, cmd = model.SearchInput.Update(msg)
	
	if newValue := model.SearchInput.Value(); newValue != prevValue {
		if newValue == "" {
			model.ResetSearch()
		} else {
			model.FilterChoices()
		}
	}
	return model, cmd
}

func (model *SimpleChooseModel) FilterChoices() {
	searchTerm := strings.ToLower(model.SearchInput.Value())
	if searchTerm == "" {
		model.FilteredChoices = model.Choices
	} else {
		filtered := make([]string, 0, len(model.Choices)/2)
		for _, choice := range model.Choices {
			if strings.Contains(strings.ToLower(choice), searchTerm) {
				filtered = append(filtered, choice)
			}
		}
		model.FilteredChoices = filtered
	}
	model.Cursor = 0
	model.ViewportStart = 0
}


func (model *SimpleChooseModel) View() string {
	var s strings.Builder
	s.Grow(1024)
	
	s.WriteString(Styles.Title.Render(model.Prompt))
	s.WriteString(" ")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.Secondary)).Render("[type to search]"))
	s.WriteString(": ")
	s.WriteString(model.SearchInput.Value())
	s.WriteString("\n")

	choiceCount := len(model.FilteredChoices)
	if choiceCount == 0 {
		s.WriteString("  No matches found\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.Muted)).Render("(esc to clear search, ctrl + c to quit)"))
		return s.String()
	}

	viewportEnd := model.ViewportStart + model.MaxVisible
	if viewportEnd > choiceCount {
		viewportEnd = choiceCount
	}

	if model.ViewportStart > 0 {
		s.WriteString(Styles.Status(Colors.Muted).Render("    ↑ more above"))
		s.WriteString("\n")
	}

	for i := model.ViewportStart; i < viewportEnd; i++ {
		if i == model.Cursor {
			s.WriteString(Styles.Selected.Render("▶ " + model.FilteredChoices[i]))
		} else {
			s.WriteString(Styles.Item.Render("  " + model.FilteredChoices[i]))
		}
		s.WriteString("\n")
	}

	if viewportEnd < choiceCount {
		s.WriteString(Styles.Status(Colors.Muted).Render("    ↓ more below"))
		s.WriteString("\n")
	}

	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.Muted)).Render("(↑/↓ Navigate, enter to select, esc to clear search, ctrl + c to quit)"))
	return s.String()
}

func (model MultiSelectModel) Init() tea.Cmd {
	return nil
}

func (model MultiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		switch messageLocal.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "up", "k":
			if model.Cursor > 0 {
				model.Cursor--
			}
		case "down", "j":
			if model.Cursor < len(model.Choices)-1 {
				model.Cursor++
			}
		case " ":
			if model.Selected[model.Cursor] {
				delete(model.Selected, model.Cursor)
			} else {
				model.Selected[model.Cursor] = true
			}
		case "enter":
			return model, tea.Quit
		}
	}
	return model, nil
}

func (model MultiSelectModel) View() string {
	s := Styles.Title.Render(model.Prompt) + "\n"
	s += Styles.Status(Colors.Info).Render("Use arrow keys to Navigate, space to select, enter to confirm") + "\n"
	for i, choice := range model.Choices {
		cursor := " "
		if model.Cursor == i {
			cursor = ">"
		}
		checked := " "
		if model.Selected[i] {
			checked = "✓"
		}
		line := cursor + " [" + checked + "] " + choice
		if model.Cursor == i {
			s += Styles.Selected.Render(line) + "\n"
		} else {
			s += Styles.Item.Render(line) + "\n"
		}
	}
	return s
}

func (model InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		switch messageLocal.Type {
		case tea.KeyEnter:
			return model, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			model.Cancelled = true
			return model, tea.Quit
		}
	}

	model.TextInput, cmd = model.TextInput.Update(msg)
	return model, cmd
}

func (model InputModel) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n%s",
		Styles.Title.Render(model.Prompt),
		model.TextInput.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color(Colors.Muted)).Render("(esc to quit)"))
}

func (model ConfirmModel) Init() tea.Cmd {
	return nil
}

func (model ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		switch messageLocal.String() {
		case "y", "Y":
			model.Confirmed = true
			return model, tea.Quit
		case "n", "N":
			model.Confirmed = false
			return model, tea.Quit
		case "ctrl+c", "esc":
			model.Confirmed = false
			return model, tea.Quit
		case "enter":
			model.Confirmed = model.DefaultValue
			return model, tea.Quit
		}
	}
	return model, nil
}

func (model ConfirmModel) View() string {
	defaultHint := "n"
	if model.DefaultValue {
		defaultHint = "Y"
	}
	return "\n" + Styles.Title.Render(model.Prompt) + "\n(Y/n) [default: " + defaultHint + "]"
}
