package notifiers

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// New creates a notifier.
func New() *Notifier {
	return &Notifier{
		MessageLogger: log.New(os.Stdout, "[message]: ", log.Ldate|log.Ltime),
		ErrorLogger:   log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
	}
}

// SendMessage sends a message to the notifier.
func (notifier *Notifier) SendMessage(val string) *Notifier {
	notifier.MessageLogger.Println(val)
	return notifier
}

// SendMessageAndTrace sends an error to the notifier and traces the runtime caller.
func (notifier *Notifier) SendMessageAndTrace(val string, skip int) *Notifier {
	_, file, line, ok := runtime.Caller(skip + 1)
	if ok {
		notifier.MessageLogger.Println(fmt.Sprintf("%s:%d %s", file, line, val))
	} else {
		notifier.MessageLogger.Println(val)
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
