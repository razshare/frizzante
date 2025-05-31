package frizzante

import (
	"fmt"
	"os"
)

type Notifier struct {
	messageFile *os.File
	errorFile   *os.File
}

// NewNotifier creates a notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		messageFile: os.Stdout,
		errorFile:   os.Stderr,
	}
}

func (notifier *Notifier) WithMessageFile(file *os.File) *Notifier {
	notifier.messageFile = file
	return notifier
}

func (notifier *Notifier) WithErrorFile(file *os.File) *Notifier {
	notifier.errorFile = file
	return notifier
}

// SendMessage sends a message to the notifier.
func (notifier *Notifier) SendMessage(message string) *Notifier {
	_, errorLocal := notifier.messageFile.WriteString(message + "\n")
	if errorLocal != nil {
		fmt.Printf("notifier could not write to message file")
	}
	return notifier
}

// SendError sends an error to the notifier.
// This will not limit the size of said messages.
func (notifier *Notifier) SendError(err error) *Notifier {
	_, errorLocal := notifier.errorFile.WriteString(err.Error() + "\n")
	if errorLocal != nil {
		fmt.Printf("notifier could not write to error file")
	}
	return notifier
}
