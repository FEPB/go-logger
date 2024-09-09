package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log       *zap.Logger
	atomLevel zap.AtomicLevel
)

// Configures the default logger for the application
func init() {
	atomLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

	config := zap.Config{
		Level:       atomLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	baseLogger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	log = baseLogger
	defer log.Sync()
}

// encoderConfig defines the default encoding configuration for the log
func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Fatal logs a message at FatalLevel then calls os.Exit(1). The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// WithFields returns a new logger with the provided fields.
func WithFields(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// SetLogLevel allows changing the log level dynamically.
func SetLogLevel(level zapcore.Level) {
	atomLevel.SetLevel(level)
}
