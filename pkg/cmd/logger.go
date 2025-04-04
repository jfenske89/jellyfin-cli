package cmd

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// initLogger initializes and returns a zap logger
func initLogger() (*zap.SugaredLogger, zap.AtomicLevel) {
	// Create an atomic level that can be changed at runtime
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.InfoLevel) // Default level

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	// Create a console encoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create a core that writes to stderr
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), atom),
	)

	// Create a logger with the core
	logger := zap.New(core)

	return logger.Sugar(), atom
}

// updateLogLevel updates the log level based on the provided level string
func updateLogLevel(levelStr string) {
	// Get the global logger
	logger := zap.L().Sugar()

	// Parse the level string
	level, err := zap.ParseAtomicLevel(strings.ToUpper(levelStr))
	if err != nil {
		logger.Warnw("Invalid log level, using INFO", "level", levelStr, "error", err)
		return
	}

	// Update the global logger's level
	zap.L().Core().Enabled(level.Level())
}
