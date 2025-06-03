package frz

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

// SendMessageAndTrace sends an error to the notifier and traces the runtime caller.
func (notifier *Notifier) SendMessageAndTrace(message string, skip int) *Notifier {
	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		notifier.messageLogger.Println(fmt.Sprintf("%s:%d %s", file, line, message))
	} else {
		notifier.messageLogger.Println(message)
	}
	return notifier
}

// SendError sends an error to the notifier.
func (notifier *Notifier) SendError(err error) *Notifier {
	notifier.errorLogger.Println(err.Error())
	return notifier
}

// SendErrorAndTrace sends an error to the notifier and traces the runtime caller.
func (notifier *Notifier) SendErrorAndTrace(err error, skip int) *Notifier {
	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		notifier.errorLogger.Println(fmt.Sprintf("%s:%d %s", file, line, err.Error()))
	} else {
		notifier.errorLogger.Println(err.Error())
	}
	return notifier
}
