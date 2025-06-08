package frz

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

var errorBuffer = bytes.NewBufferString("")
var errorLog = log.New(errorBuffer, "[error]: ", log.Ldate|log.Ltime)

var messageBuffer = bytes.NewBufferString("")
var messageLog = log.New(messageBuffer, "[message]: ", log.Ldate|log.Ltime)

func TestSendMessage(t *testing.T) {
	notifier := NewNotifier().WithMessageLogger(messageLog)
	notifier.SendMessage("hello")
	readString, readError := messageBuffer.ReadString('\n')
	if readError != nil {
		t.Fatal(readError)
	}
	if !strings.HasPrefix(readString, "[message]:") {
		t.Fatal("notifier message filed")
		return
	}
	if !strings.HasSuffix(readString, "hello\n") {
		t.Fatal("notifier message filed")
		return
	}
}

func TestSendError(t *testing.T) {
	notifier := NewNotifier().WithErrorLogger(errorLog)
	notifier.SendError(errors.New("this is an error"))
	readString, readError := errorBuffer.ReadString('\n')
	if readError != nil {
		t.Fatal(readError)
	}
	if !strings.HasPrefix(readString, "[error]:") {
		t.Fatal("notifier error filed")
		return
	}
	if !strings.HasSuffix(readString, "this is an error\n") {
		t.Fatal("notifier error filed")
		return
	}
}
