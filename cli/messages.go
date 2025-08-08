package cli

import "github.com/pterm/pterm"

// Fatalf shows a fatal message and terminates the application.
func Fatalf(template string, vars ...any) {
	pterm.Fatal.Printfln(template, vars...)
}

// Warningf shows a warning message.
func Warningf(template string, vars ...any) {
	pterm.Warning.Printfln(template, vars...)
}

// Infof shows an info message.
func Infof(template string, vars ...any) {
	pterm.Info.Printfln(template, vars...)
}

// Successf shows a success message.
func Successf(template string, vars ...any) {
	pterm.Success.Printfln(template, vars...)
}

// Warning shows a warning message.
func Warning(vars ...any) {
	pterm.Warning.Println(vars...)
}

// Fatal shows a fatal message and terminates the application.
func Fatal(vars ...any) {
	pterm.Fatal.Println(vars...)
}

// Info shows an info message.
func Info(vars ...any) {
	pterm.Info.Println(vars...)
}

// Success shows a success message.
func Success(vars ...any) {
	pterm.Success.Println(vars...)
}

// Section shows the name of a section using Markdown semantics.
func Section(vars ...any) {
	pterm.DefaultSection.Println(vars...)
}
