package config

import (
	"errors"
	"github.com/Direct-Debit/go-commons/errlib"
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

func GetDef(key string, def interface{}) interface{} {
	lock.Lock()
	defer lock.Unlock()

	once.Do(func() {
		var err error
		conf, err = toml.LoadFile("config.toml")
		errlib.WarnError(err, "Error loading config")
	})
	if conf == (*toml.Tree)(nil) {
		return def
	}

	val := conf.Get(key)
	if val == nil {
		return def
	}
	return val
}

func Get(key string) interface{} {
	v := GetDef(key, nil)
	if v == nil {
		errlib.FatalError(errors.New(key), "Config variable missing")
	}
	return v
}

func GetStrDef(key string, def string) string {
	return GetDef(key, def).(string)
}

func GetStr(key string) string {
	return Get(key).(string)
}

func GetBoolDef(key string, def bool) bool {
	return GetDef(key, def).(bool)
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
	levelStr := GetStrDef("log_level", "DEBUG")
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
		if GetBoolDef("debug", true) {
			level = log.DebugLevel
		}
		log.Warnf("CONFIG WARNING: Could not parse log_level %v, defaulting to %v level", levelStr, level)
		return level
	}
}
