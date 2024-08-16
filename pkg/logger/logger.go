package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wojciechpawlinow/usermanagement/internal/config"
)

var l *zap.SugaredLogger

// Setup initializes the logging infrastructure based on the provided configuration.
// It should be called once during the startup of the application.
func Setup(cfg config.Provider) {
	var level zapcore.Level

	switch cfg.GetString("log_level") {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		level,
	)

	l = zap.New(core).Sugar()
}

// Debug logs a debug message with the given fields.
func Debug(args ...interface{}) {
	l.Debug(args...)
}

// Info logs an informational message with the given fields.
func Info(args ...interface{}) {
	l.Info(args...)
}

// Error logs an error message with the given fields.
func Error(args ...interface{}) {
	l.Error(args...)
}

// Fatal logs a fatal message, then exits the application.
func Fatal(args ...interface{}) {
	l.Fatal(args...)
}
