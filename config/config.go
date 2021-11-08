package config

import (
	"github.com/Direct-Debit/go-commons/config/tomlold"
	"github.com/Direct-Debit/go-commons/errlib"
	"github.com/Direct-Debit/go-commons/stdext"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Provider interface {
	GetDef(key string, def interface{}) (interface{}, error)
	Get(key string) (interface{}, error)
}

var conf Provider = tomlold.NewReader()

func SetProvider(c Provider) { conf = c }

func GetProvider() Provider { return conf }

func GetDef(key string, def interface{}) interface{} {
	val, err := conf.GetDef(key, def)
	errlib.DebugError(err, "Error reading %s config setting (defaulting to %v)", key, def)
	return val
}

func Get(key string) interface{} {
	v, err := conf.Get(key)
	errlib.FatalError(err, "Error reading %s config setting", key)
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

func GetInt64(key string) int64 {
	return Get(key).(int64)
}

func GetInt64Def(key string, def int64) int64 {
	return GetDef(key, def).(int64)
}

func GetInt(key string) int {
	return int(GetInt64(key))
}

func GetIntDef(key string, def int) int {
	return int(GetInt64Def(key, int64(def)))
}

func GetFloat(key string) float64 {
	return Get(key).(float64)
}

func GetStrListDef(key string, def []string) []string {
	val := GetDef(key, nil)
	if val == nil {
		return def
	}

	return stdext.SliceInterfaceToString(val.([]interface{}))
}

func GetStrList(key string) []string {
	list := Get(key).([]interface{})
	return stdext.SliceInterfaceToString(list)
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
	case "PANIC":
		return log.PanicLevel
	default:
		level := log.InfoLevel
		if GetBoolDef("debug", true) {
			level = log.DebugLevel
		}
		log.Warnf("CONFIG WARNING: Could not parse log_level %v, defaulting to %v level", levelStr, level)
		return level
	}
}
