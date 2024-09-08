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
func Info(msg string, fields ...interface{}) {
	log.Info(msg, convertToZapFields(fields...)...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...interface{}) {
	log.Error(msg, convertToZapFields(fields...)...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...interface{}) {
	log.Debug(msg, convertToZapFields(fields...)...)
}

// Fatal logs a message at FatalLevel then calls os.Exit(1). The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func Fatal(msg string, fields ...interface{}) {
	log.Fatal(msg, convertToZapFields(fields...)...)
}

// convertToZapFields converts a dynamic list of interface{} into zap.Fields based on their type
func convertToZapFields(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)/2)

	// Expecting pairs of key-value like: key, value, key, value
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			// Skip if the key is not a string
			continue
		}
		value := fields[i+1]

		// Handle different types for value
		switch v := value.(type) {
		case string:
			zapFields = append(zapFields, zap.String(key, v))
		case int:
			zapFields = append(zapFields, zap.Int(key, v))
		case bool:
			zapFields = append(zapFields, zap.Bool(key, v))
		case float64:
			zapFields = append(zapFields, zap.Float64(key, v))
		default:
			// Fallback for other types using reflection
			zapFields = append(zapFields, zap.Reflect(key, v))
		}
	}
	return zapFields
}

// WithFields returns a new logger with the provided fields.
func WithFields(fields ...interface{}) *zap.Logger {
	return log.With(convertToZapFields(fields...)...)
}

// SetLogLevel allows changing the log level dynamically.
func SetLogLevel(level zapcore.Level) {
	atomLevel.SetLevel(level)
}
