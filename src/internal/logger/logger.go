package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel int8) (*zap.Logger, error) {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(logLevel))
	return loggerCfg.Build()
}
