package filesystem

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

func mockLogger(t *testing.T) *zap.Logger {
	t.Helper()

	zapCore, _ := observer.New(zapcore.DebugLevel)
	return zap.New(zapCore)
}

func TestPlatformNamer_GetPlatformNameForAtcoCode(t *testing.T) {
	t.Run(`Given an AtcoCode has an associated platform name
When GetPlatformNameForAtcoCode is called with the AtcoCode
Then the name of the platform is returned`, func(t *testing.T) {
		// Given
		logger := mockLogger(t)

		platformNamer := NewPlatformNamer(logger)

		// When
		platformName, err := platformNamer.GetPlatformNameForAtcoCode("9400ZZMASTP3")

		// Then
		assert.Nil(t, err)
		assert.NotNil(t, platformName)
		assert.Equal(t, "B", *platformName)
	})

	t.Run(`Given an AtcoCode does not have an associated platform name
When GetPlatformNameForAtcoCode is called with the AtcoCode
Then nil is returned`, func(t *testing.T) {
		// Given
		logger := mockLogger(t)

		platformNamer := NewPlatformNamer(logger)

		// When
		platformName, err := platformNamer.GetPlatformNameForAtcoCode("9400ZZMAOLD1")

		// Then
		assert.Nil(t, err)
		assert.Nil(t, platformName)
	})
}
