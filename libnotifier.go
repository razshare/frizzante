package frizzante

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type Notifier struct {
	messageLogger *log.Logger
	errorLogger   *log.Logger
}

// NewNotifier creates a notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		messageLogger: log.New(os.Stdout, "[message]: ", log.Ldate|log.Ltime),
		errorLogger:   log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
	}
}

func (notifier *Notifier) WithMessageLogger(logger *log.Logger) *Notifier {
	notifier.messageLogger = logger
	return notifier
}

func (notifier *Notifier) WithErrorLogger(logger *log.Logger) *Notifier {
	notifier.messageLogger = logger
	return notifier
}

// SendMessage sends a message to the notifier.
func (notifier *Notifier) SendMessage(message string) *Notifier {
	notifier.messageLogger.Println(message)
	return notifier
}

// SendErrorAndTrace sends an error to the notifier.
func (notifier *Notifier) SendErrorAndTrace(err error) *Notifier {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		notifier.messageLogger.Println(fmt.Sprintf("%s:%d %s", file, line, err.Error()))
	} else {
		notifier.messageLogger.Println(err.Error())
	}
	return notifier
}

// SendError sends an error to the notifier without tracing the caller.
func (notifier *Notifier) SendError(err error) *Notifier {
	notifier.messageLogger.Println(err.Error())
	return notifier
}
