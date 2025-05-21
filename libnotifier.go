package frizzante

import (
	"fmt"
	"os"
)

type Notifier struct {
	errorFile   *os.File
	messageFile *os.File
}

// NewNotifier creates a notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		errorFile:   os.Stderr,
		messageFile: os.Stdout,
	}
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

// SendMessage sends a message to the notifier.
func (notifier *Notifier) SendMessage(message string) *Notifier {
	_, errorLocal := notifier.messageFile.WriteString(message + "\n")
	if errorLocal != nil {
		fmt.Printf("notifier could not write to message file")
	}
	return notifier
}
