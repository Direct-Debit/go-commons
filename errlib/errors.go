package errlib

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func Error(err error, message string, printer func(...interface{})) bool {
	if err != nil {
		printer(fmt.Sprintf("%v: %v", message, err))
		return true
	}
	return false
}

func FatalError(err error, message string) {
	Error(err, message, log.Fatal)
}

func PanicError(err error, message string) {
	Error(err, message, log.Panic)
}

// If err != nil, print message: err to the log and return true
func ErrorError(err error, message string) bool {
	return Error(err, message, log.Error)
}

func WarnError(err error, message string) bool {
	return Error(err, message, log.Warn)
}

func DebugError(err error, message string) bool {
	return Error(err, message, log.Debug)
}
