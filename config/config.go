package config

import (
	"errors"
	"fmt"
	"github.com/Direct-Debit/gocommons/errlib"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
)

var (
	conf *toml.Tree
	once sync.Once
	lock sync.Mutex
)

func Get(key string) interface{} {
	lock.Lock()
	defer lock.Unlock()

	once.Do(func() {
		var err error
		conf, err = toml.LoadFile("config.toml")
		errlib.FatalError(err, "Error loading config")
	})

	val := conf.Get(key)
	if val == nil {
		errlib.FatalError(errors.New(key), "Config variable missing")
	}
	return val
}

func GetStr(key string) string {
	return Get(key).(string)
}

func GetBool(key string) bool {
	return Get(key).(bool)
}

func GetInt(key string) int {
	return Get(key).(int)
}

func GetFloat(key string) float64 {
	return Get(key).(float64)
}

func GetStrList(key string) []string {
	list := conf.Get(key).([]interface{})
	result := make([]string, len(list))
	for i, l := range list {
		result[i] = l.(string)
	}
	return result
}

func GetLogLevel() log.Level {
	levelStr := GetStr("log_level")
	switch strings.ToUpper(levelStr) {
	case "TRACE":
		return log.TraceLevel
	case "DEBUG":
		return log.DebugLevel
	case "INFO":
		return log.InfoLevel
	case "WARNING":
		return log.WarnLevel
	case "ERROR":
		return log.ErrorLevel
	case "ERR":
		return log.ErrorLevel
	case "FATAL":
		return log.FatalLevel
	default:
		level := log.InfoLevel
		if GetBool("debug") {
			level = log.DebugLevel
		}
		log.Error(
			fmt.Sprintf("CONFIG ERROR: Could not parse log_level %v, defaulting to %v level", levelStr, level))
		return level
	}
}
