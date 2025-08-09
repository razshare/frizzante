package cli

import "github.com/razshare/frizzante/tui/messages"

// Fatalf shows a fatal message and terminates the application.
func Fatalf(format string, vars ...any) {
	messages.Fatalf(format, vars...)
}

// Warningf shows a warning message.
func Warningf(format string, vars ...any) {
	messages.Warningf(format, vars...)
}

// Infof shows an info message.
func Infof(format string, vars ...any) {
	messages.Infof(format, vars...)
}

// Successf shows a success message.
func Successf(format string, vars ...any) {
	messages.Successf(format, vars...)
}

// Warning shows a warning message.
func Warning(vars ...any) {
	messages.Warning(vars...)
}

// Fatal shows a fatal message and terminates the application.
func Fatal(vars ...any) {
	messages.Fatal(vars...)
}

// Info shows an info message.
func Info(vars ...any) {
	messages.Info(vars...)
}

// Success shows a success message.
func Success(vars ...any) {
	messages.Success(vars...)
}

// Section shows the name of a section using Markdown semantics.
func Section(vars ...any) {
	messages.Section(vars...)
}
