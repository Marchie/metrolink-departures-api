package sqs_test

import (
	"context"
	mock_core "github.com/Marchie/tf-experiment/lambda/internal/mocks/core"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/sqs"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
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

func TestMetrolinkDeparturesDataLoader_Handler(t *testing.T) {
	t.Run(`Given Metrolink departures are loaded
When Handler is called
Then no error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		metrolinkDeparturesLoader := mock_core.NewMockMetrolinkDeparturesLoader(ctrl)
		metrolinkDeparturesLoader.EXPECT().Load(ctx).Return(nil)

		metrolinkDeparturesDataLoader := sqs.NewMetrolinkDeparturesDataLoader(logger, metrolinkDeparturesLoader)

		// When
		err := metrolinkDeparturesDataLoader.Handler(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given an error occurs loading Metrolink departures
When Handler is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		metrolinkDeparturesLoader := mock_core.NewMockMetrolinkDeparturesLoader(ctrl)
		metrolinkDeparturesLoaderErr := errors.New("FUBAR")
		metrolinkDeparturesLoader.EXPECT().Load(ctx).Return(metrolinkDeparturesLoaderErr)

		metrolinkDeparturesDataLoader := sqs.NewMetrolinkDeparturesDataLoader(logger, metrolinkDeparturesLoader)

		// When
		err := metrolinkDeparturesDataLoader.Handler(ctx)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, metrolinkDeparturesLoaderErr, err)
	})
}
