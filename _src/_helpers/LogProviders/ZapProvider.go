package logproviders

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZapLogger(isDevelopment bool, logFilePath string, debugLevel zapcore.Level, initialFields map[string]any) *zap.Logger {
	// Configure Zap logger options
	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
	} else {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.NameKey = "name"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.LineEnding = "\n"
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 time format
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.DisableCaller = false
	// Customize logging level for development
	config.Level.SetLevel(debugLevel)
	config.Development = isDevelopment
	config.OutputPaths = []string{"stdout", logFilePath}
	config.ErrorOutputPaths = []string{"stderr"}
	config.DisableStacktrace = false
	config.InitialFields = initialFields
	// Build the logger
	logger, err := config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync() // Flushes buffer, if any
	return logger
}
