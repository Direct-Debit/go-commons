package errlib

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func Error(printer func(...interface{}), err error, format string, a ...interface{}) bool {
	if err != nil {
		msg := fmt.Sprintf(format, a)
		printer(fmt.Sprintf("%v: %v", msg, err))
		return true
	}
	return false
}

func FatalError(err error, format string, a ...interface{}) {
	Error(log.Fatal, err, format, a)
}

func PanicError(err error, format string, a ...interface{}) {
	Error(log.Panic, err, format, a)
}

// ErrorError checks if err != nil, in which case it logs the error with fmt.Sprintf(format, a) prepended
func ErrorError(err error, format string, a ...interface{}) bool {
	return Error(log.Error, err, format, a)
}

func WarnError(err error, format string, a ...interface{}) bool {
	return Error(log.Warn, err, format, a)
}

func InfoError(err error, format string, a ...interface{}) bool {
	return Error(log.Info, err, format, a)
}

func DebugError(err error, format string, a ...interface{}) bool {
	return Error(log.Debug, err, format, a)
}
