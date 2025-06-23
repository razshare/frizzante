package libnotifier

import "log"

type Notifier struct {
	MessageLogger *log.Logger
	ErrorLogger   *log.Logger
}
