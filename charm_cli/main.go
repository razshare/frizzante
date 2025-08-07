package main

import (
	"main/cli"
	//"time"
)

func main() {
	// Test status messages (replacing pterm.Info, Success, Warning, Fatal)
	// cli.CharmSection("Charm CLI - Replacing pterm Functions")
	// cli.CharmInfo("Testing charm_cli as replacement for pterm")
	// cli.CharmSuccess("Success messages (replaces pterm.Success)")
	// cli.CharmWarning("Warning messages (replaces pterm.Warning)")
	// cli.CharmError("Error messages (replaces pterm.Fatal)")
	
	// Test interactive select (replacing pterm.DefaultInteractiveSelect)
	cli.CharmSection("Interactive Menu (replaces pterm.DefaultInteractiveSelect)")
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
	menuChoice, err := cli.CharmChoose("Pick a Frizzante option", frizzanteOptions)
	if err != nil {
		cli.CharmError("Menu selection failed: " + err.Error())
	} else {
		cli.CharmSuccess("You selected: " + menuChoice)
	}
	
	// // Test text input (replacing pterm.DefaultInteractiveTextInput)
	// cli.CharmSection("Text Input (replaces pterm.DefaultInteractiveTextInput)")
	// projectName, err := cli.CharmInput("Enter project name:")
	// if err != nil {
	// 	cli.CharmError("Text input failed: " + err.Error())
	// } else {
	// 	cli.CharmSuccess("Project name: " + projectName)
	// }
	
	// // Test confirm (replacing pterm.DefaultInteractiveConfirm)
	// cli.CharmSection("Confirmation (replaces pterm.DefaultInteractiveConfirm)")
	// overwrite, err := cli.CharmConfirm("Feature already exists, overwrite?", false)
	// if err != nil {
	// 	cli.CharmError("Confirm failed: " + err.Error())
	// } else if overwrite {
	// 	cli.CharmSuccess("Will overwrite existing feature")
	// } else {
	// 	cli.CharmInfo("Skipping feature installation")
	// }
	
	// // Test multi-select (replacing pterm.DefaultInteractiveMultiselect)
	// cli.CharmSection("Multi-Select (replaces pterm.DefaultInteractiveMultiselect)")
	// features := []string{"Core", "Forms", "Links", "Air", "Bun", "Sqlc"}
	// selectedFeatures, err := cli.CharmMultiSelect("Pick features to add:", features)
	// if err != nil {
	// 	cli.CharmError("Multi-select failed: " + err.Error())
	// } else {
	// 	if len(selectedFeatures) > 0 {
	// 		cli.CharmSuccess("Selected features:")
	// 		for _, feature := range selectedFeatures {
	// 			cli.CharmInfo("- " + feature)
	// 		}
	// 	} else {
	// 		cli.CharmInfo("No features selected")
	// 	}
	// }
	
	// // Test platform selection (another pterm.DefaultInteractiveSelect example)
	// cli.CharmSection("Platform Selection Demo")
	// platforms := []string{
	// 	"Linux/amd64",
	// 	"Linux/arm64", 
	// 	"Darwin/amd64",
	// 	"Darwin/arm64",
	// 	"Windows/amd64",
	// 	"Windows/arm64",
	// }
	// platform, err := cli.CharmChoose("Pick a platform:", platforms)
	// if err != nil {
	// 	cli.CharmError("Platform selection failed: " + err.Error())
	// } else {
	// 	cli.CharmSuccess("Target platform: " + platform)
	// }
	
	// // Test table (replacing pterm.DefaultTable)
	// cli.CharmSection("Table Display (replaces pterm.DefaultTable)")
	// tableHeaders := []string{"Feature Name", "Description", "Status"}
	// tableRows := [][]string{
	// 	{"Core", "Base frizzante functionality with view rendering and state management", "Required"},
	// 	{"Forms", "Enhanced form component with web standards support", "Optional"},
	// 	{"Links", "Enhanced link component with additional features", "Optional"},
	// 	{"Bun", "JavaScript runtime and package manager", "Development"},
	// 	{"Sqlc", "Type-safe SQL code generator", "Database"},
	// }
	// cli.CharmTable(tableHeaders, tableRows)
	
	// // Test spinner (replacing pterm.DefaultSpinner)
	// cli.CharmSection("Spinner (replaces pterm.DefaultSpinner)")
	// installSpinner := cli.CharmSpinner("Installing bun from GitHub releases...")
	// installSpinner.Start()
	// time.Sleep(4 * time.Second)
	// installSpinner.Stop()
	// cli.CharmSuccess("Bun installed successfully!")
	
	// // Test Docker help output
	// cli.CharmSection("Docker Help Output")
	// cli.CharmDockerHelp()
	
	// cli.CharmSection("pterm Replacement Test Complete")
	// cli.CharmSuccess("All pterm functions successfully replaced with charm_cli!")
	// cli.CharmInfo("Ready to replace pterm imports in the main codebase")
}
