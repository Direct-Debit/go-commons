package config

import (
	"github.com/Direct-Debit/go-commons/format"
	log "github.com/sirupsen/logrus"
)

func SetupLogs() {
	log.SetLevel(GetLogLevel())
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: format.DateTimeShortDashes})
}
