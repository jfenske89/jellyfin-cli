package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildLogger() (*zap.SugaredLogger, zap.AtomicLevel) {
	atom := zap.NewAtomicLevel()
	stdout := zapcore.AddSync(os.Stderr)

	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	encoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(zapcore.NewCore(encoder, stdout, atom))

	return zap.New(core).Sugar(), atom
}
