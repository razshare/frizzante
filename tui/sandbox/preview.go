package sandbox

import (
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"github.com/razshare/frizzante/tui/table"
	"time"
)

func Preview() {
	// Test status messages (replacing pterm.Info, Success, Warning, Fatal)
	messages.Section("Charm CLI - Replacing pterm Functions")
	messages.Info("Testing charm_cli as replacement for pterm")
	messages.Success("Success messages (replaces pterm.Success)")
	messages.Warning("Warning messages (replaces pterm.Warning)")
	messages.Fatal("Error messages (replaces pterm.Fatal)")

	// Test interactive select (replacing pterm.DefaultInteractiveSelect)
	messages.Section("Interactive Menu (replaces pterm.DefaultInteractiveSelect)")
	frizzanteOptions := []string{
		"Help",
		"Version",
		"Create Project",
		"Add Features",
		"Test",
		"Package",
		"Dev",
		"Build",
		"Configure",
	}
	menuChoice, err := singleselect.Send("Pick a Frizzante option", frizzanteOptions)
	if err != nil {
		messages.Fatal("Menu selection failed: " + err.Error())
	} else {
		messages.Success("You selected: " + menuChoice)
	}

	// // Test text input (replacing pterm.DefaultInteractiveTextInput)
	messages.Section("Text Input (replaces pterm.DefaultInteractiveTextInput)")
	projectName, err := input.Send("Enter project name:")
	if err != nil {
		messages.Fatal("Text input failed: " + err.Error())
	} else {
		messages.Success("Project name: " + projectName)
	}

	// // Test confirm (replacing pterm.DefaultInteractiveConfirm)
	messages.Section("Confirmation (replaces pterm.DefaultInteractiveConfirm)")
	overwrite, err := confirm.Send("Feature already exists, overwrite?", false)
	if err != nil {
		messages.Fatal("Confirm failed: " + err.Error())
	} else if overwrite {
		messages.Success("Will overwrite existing feature")
	} else {
		messages.Info("Skipping feature installation")
	}

	// // Test multi-select (replacing pterm.DefaultInteractiveMultiselect)
	messages.Section("Multi-Select (replaces pterm.DefaultInteractiveMultiselect)")
	features := []string{"Core", "Forms", "Links", "Air", "Bun", "Sqlc"}
	selectedFeatures, err := multiselect.Send("Pick features to add:", features)
	if err != nil {
		messages.Fatal("Multi-select failed: " + err.Error())
	} else {
		if len(selectedFeatures) > 0 {
			messages.Success("Selected features:")
			for _, feature := range selectedFeatures {
				messages.Info("- " + feature)
			}
		} else {
			messages.Info("No features selected")
		}
	}

	// // Test platform selection (another pterm.DefaultInteractiveSelect example)
	messages.Section("Platform Selection Demo")
	platforms := []string{
		"Linux/amd64",
		"Linux/arm64",
		"Darwin/amd64",
		"Darwin/arm64",
		"Windows/amd64",
		"Windows/arm64",
	}
	platform, err := singleselect.Send("Pick a platform:", platforms)
	if err != nil {
		messages.Fatal("Platform selection failed: " + err.Error())
	} else {
		messages.Success("Target platform: " + platform)
	}

	// // Test table (replacing pterm.DefaultTable)
	messages.Section("Table Display (replaces pterm.DefaultTable)")
	tableHeaders := []string{"Feature Name", "Description", "Status"}
	tableRows := [][]string{
		{"Core", "Base frizzante functionality with view rendering and state management", "Required"},
		{"Forms", "Enhanced form component with web standards support", "Optional"},
		{"Links", "Enhanced link component with additional features", "Optional"},
		{"Bun", "JavaScript runtime and package manager", "Development"},
		{"Sqlc", "Type-safe SQL code generator", "Database"},
	}
	table.Send(tableHeaders, tableRows)

	// // Test spinner (replacing pterm.DefaultSpinner)
	messages.Section("Spinner (replaces pterm.DefaultSpinner)")
	installSpinner := spinner.New("Installing bun from GitHub releases...")
	spinner.Start(installSpinner)
	time.Sleep(4 * time.Second)
	spinner.Stop(installSpinner)
	messages.Success("Bun installed successfully!")

	// // Test Docker help output
	messages.Section("Docker Help Output")
	messages.DockerHelp()

	messages.Section("pterm Replacement Test Complete")
	messages.Success("All pterm functions successfully replaced with charm_cli!")
	messages.Info("Ready to replace pterm imports in the main codebase")
}
