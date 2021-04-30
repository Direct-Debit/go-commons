package errlib

import (
	"fmt"
)

func Error(printer func(string), err error, format string, a ...interface{}) bool {
	if err != nil {
		msg := fmt.Sprintf(format, a...)
		printer(fmt.Sprintf("%v: %v", msg, err))
		return true
	}
	return false
}

func FatalError(err error, format string, a ...interface{}) {
	Error(fatalFunc, err, format, a)
}

func PanicError(err error, format string, a ...interface{}) {
	Error(panicFunc, err, format, a)
}

// ErrorError checks if err != nil, in which case it logs the error with fmt.Sprintf(format, a) prepended
func ErrorError(err error, format string, a ...interface{}) bool {
	return Error(errorFunc, err, format, a)
}

func WarnError(err error, format string, a ...interface{}) bool {
	return Error(warnFunc, err, format, a)
}

func InfoError(err error, format string, a ...interface{}) bool {
	return Error(infoFunc, err, format, a)
}

func DebugError(err error, format string, a ...interface{}) bool {
	return Error(debugFunc, err, format, a)
}

func TraceError(err error, format string, a ...interface{}) bool {
	return Error(traceFunc, err, format, a)
}
