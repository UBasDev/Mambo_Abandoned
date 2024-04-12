package logproviders

import (
	"github.com/pkg/errors"
)

func LoadLogger(lc LogConfig) error {
	loggerType := lc.Code
	err := GetLogFactoryBuilder(loggerType).Build(&lc)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

var logfactoryBuilderMap = map[string]logFbInterface{
	"zap": &ZapFactory{},
	//Buraya başka log factoryleri de ekleyebiliriz
}

// interface for logger factory
type logFbInterface interface {
	Build(*LogConfig) error
}

// accessors for factoryBuilderMap
func GetLogFactoryBuilder(key string) logFbInterface {
	return logfactoryBuilderMap[key]
}

var Log Logger

// Logger represent common interface for logging function
type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
}

func SetLogger(newLogger Logger) {
	Log = newLogger
}

type LogConfig struct {
	// log library name
	Code string `yaml:"code"`
	// log level
	Level string `yaml:"level"`
	// show caller in log message
	EnableCaller bool `yaml:"enableCaller"`
}
