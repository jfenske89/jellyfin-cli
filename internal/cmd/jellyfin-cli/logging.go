package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildLogger() (*zap.SugaredLogger, zap.AtomicLevel) {
	atom := zap.NewAtomicLevel()
	stdout := zapcore.AddSync(os.Stderr)

	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(zapcore.NewCore(encoder, stdout, atom))

	return zap.New(core).Sugar(), atom
}
