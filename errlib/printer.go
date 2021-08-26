package errlib

import log "github.com/sirupsen/logrus"

type Logger interface {
	Fatal(string)
	Panic(string)
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
	Trace(string)
}

var loggers []Logger

func AddLogger(l Logger) {
	loggers = append(loggers, l)
}

func ClearLoggers() {
	loggers = []Logger{}
}

func fatalFunc(msg string) {
	for _, l := range loggers {
		l.Fatal(msg)
	}
	log.StandardLogger().Fatal(msg)
}

func panicFunc(msg string) {
	for _, l := range loggers {
		l.Panic(msg)
	}
	log.StandardLogger().Panic(msg)
}

func errorFunc(msg string) {
	for _, l := range loggers {
		l.Error(msg)
	}
	log.StandardLogger().Error(msg)
}

func warnFunc(msg string) {
	for _, l := range loggers {
		l.Warn(msg)
	}
	log.StandardLogger().Warn(msg)
}

func infoFunc(msg string) {
	for _, l := range loggers {
		l.Info(msg)
	}
	log.StandardLogger().Info(msg)
}

func debugFunc(msg string) {
	for _, l := range loggers {
		l.Debug(msg)
	}
	log.StandardLogger().Debug(msg)
}

func traceFunc(msg string) {
	for _, l := range loggers {
		l.Trace(msg)
	}
	log.StandardLogger().Trace(msg)
}
