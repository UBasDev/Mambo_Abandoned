package logproviders

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ZapFactory struct{}

// build zap logger
func (mf *ZapFactory) Build(lc *LogConfig) error {
	err := RegisterZapLog(*lc)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func RegisterZapLog(lc LogConfig) error {
	zLogger, err := initZapLog(lc)
	if err != nil {
		return errors.Wrap(err, "RegisterLog")
	}
	defer zLogger.Sync()
	zSugarlog := zLogger.Sugar()
	//zSugarlog.Info()

	//This is for loggerWrapper implementation
	//appLogger.SetLogger(&loggerWrapper{zaplog})

	SetLogger(zSugarlog)
	return nil
}

// initLog create logger
func initZapLog(lc LogConfig) (zap.Logger, error) {
	//Alternatif levellar: debug, info, warn, error, dpanic, panic, fatal, _min, _max
	rawJSON := []byte(`{
	 "level": "info",
     "Development": true,
     "DisableCaller": false,
	 "encoding": "console",
	 "outputPaths": ["stdout", "./logs/demo.log"],
	 "errorOutputPaths": ["stderr"],
	 "encoderConfig": {
		"timeKey":        "ts",
		"levelKey":       "level",
		"messageKey":     "msg",
        "nameKey":        "name",
		"stacktraceKey":  "stacktrace",
        "callerKey":      "caller",
		"lineEnding":     "\n\t",
        "timeEncoder":     "time",
		"levelEncoder":    "lowercaseLevel",
        "durationEncoder": "stringDuration",
        "callerEncoder":   "shortCaller"
	 }
	}`)

	var cfg zap.Config
	var zLogger *zap.Logger
	//standard configuration
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return *zLogger, errors.Wrap(err, "Unmarshal")
	}
	//customize it from configuration file
	err := customizeZapLogFromConfig(&cfg, lc)
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}
	zLogger, err = cfg.Build()
	if err != nil {
		return *zLogger, errors.Wrap(err, "cfg.Build()")
	}

	zLogger.Debug("logger construction succeeded")
	return *zLogger, nil
}

// customizeLogFromConfig customize log based on parameters from configuration file
func customizeZapLogFromConfig(cfg *zap.Config, lc LogConfig) error {
	cfg.DisableCaller = !lc.EnableCaller

	// set log level
	l := zap.NewAtomicLevel().Level()
	err := l.Set(lc.Level)
	if err != nil {
		return errors.Wrap(err, "")
	}
	cfg.Level.SetLevel(l)

	return nil
}
