package libnotifier

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// NewNotifier creates a notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		MessageLogger: log.New(os.Stdout, "[message]: ", log.Ldate|log.Ltime),
		ErrorLogger:   log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
	}
}

func (notifier *Notifier) WithMessageLogger(logger *log.Logger) *Notifier {
	notifier.MessageLogger = logger
	return notifier
}

func (notifier *Notifier) WithErrorLogger(logger *log.Logger) *Notifier {
	notifier.ErrorLogger = logger
	return notifier
}

// SendMessage sends a message to the notifier.
func (notifier *Notifier) SendMessage(message string) *Notifier {
	notifier.MessageLogger.Println(message)
	return notifier
}

// SendMessageAndTrace sends an error to the notifier and traces the runtime caller.
func (notifier *Notifier) SendMessageAndTrace(message string, skip int) *Notifier {
	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		notifier.MessageLogger.Println(fmt.Sprintf("%s:%d %s", file, line, message))
	} else {
		notifier.MessageLogger.Println(message)
	}
	return notifier
}

// SendError sends an error to the notifier.
func (notifier *Notifier) SendError(err error) *Notifier {
	notifier.ErrorLogger.Println(err.Error())
	return notifier
}

// SendErrorAndTrace sends an error to the notifier and traces the runtime caller.
func (notifier *Notifier) SendErrorAndTrace(err error, skip int) *Notifier {
	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		notifier.ErrorLogger.Println(fmt.Sprintf("%s:%d %s", file, line, err.Error()))
	} else {
		notifier.ErrorLogger.Println(err.Error())
	}
	return notifier
}
