package logger_test

import (
	"github.com/Marchie/tf-experiment/lambda/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Run(`Given a log level of zapcore.InfoLevel
When NewLogger is called
Then a configured Zap logger is returned`, func(t *testing.T) {
		// Given
		logLevel := zapcore.InfoLevel

		// When
		l, err := logger.NewLogger(int8(logLevel))

		// Then
		assert.Nil(t, err)

		assert.False(t, l.Core().Enabled(zapcore.DebugLevel))
		assert.True(t, l.Core().Enabled(zapcore.InfoLevel))
	})
}
