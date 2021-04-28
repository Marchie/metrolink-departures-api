package v1_test

import (
	"context"
	v1 "github.com/Marchie/tf-experiment/lambda/internal/core/schedule/departures/metrolink/v1"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	mock_repository "github.com/Marchie/tf-experiment/lambda/internal/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
	"time"
)

func mockLogger(t *testing.T) *zap.Logger {
	t.Helper()

	zapCore, _ := observer.New(zapcore.DebugLevel)
	return zap.New(zapCore)
}

func TestNewMetrolinkDeparturesDataLoadScheduler(t *testing.T) {
	t.Run(`Given an invalid configuration
When NewMetrolinkDeparturesDataLoadScheduler is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		horizon := time.Second * 30

		frequency := time.Duration(0)

		payload := "foo"

		now := time.Date(2021, time.March, 28, 16, 39, 19, 0, time.UTC)

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := mock_repository.NewMockEventScheduler(ctrl)

		// When
		metrolinkDeparturesDataLoadScheduler, err := v1.NewMetrolinkDeparturesDataLoadScheduler(logger, eventScheduler, horizon, frequency, payload, currentTimeFunc)

		// Then
		assert.Nil(t, metrolinkDeparturesDataLoadScheduler)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "frequency cannot be 0")
	})
}

func TestMetrolinkDeparturesDataLoadScheduler_Schedule(t *testing.T) {
	t.Run(`Given a valid configuration
When Schedule is called
Then events are returned according to the configured frequency and horizon`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		horizon := time.Second * 30

		frequency := time.Second * 8

		payload := "foo"

		now := time.Date(2021, time.March, 28, 16, 39, 19, 0, time.UTC)

		currentTimeFunc := func() time.Time {
			return now
		}

		expEvents := []*domain.Event{
			{
				StartTime: now,
				Payload:   payload,
			},
			{
				StartTime: time.Date(2021, time.March, 28, 16, 39, 27, 0, time.UTC),
				Payload:   payload,
			},
			{
				StartTime: time.Date(2021, time.March, 28, 16, 39, 35, 0, time.UTC),
				Payload:   payload,
			},
			{
				StartTime: time.Date(2021, time.March, 28, 16, 39, 43, 0, time.UTC),
				Payload:   payload,
			},
		}

		eventScheduler := mock_repository.NewMockEventScheduler(ctrl)
		eventScheduler.EXPECT().Schedule(ctx, expEvents).Return(nil)

		metrolinkDeparturesDataLoadScheduler, _ := v1.NewMetrolinkDeparturesDataLoadScheduler(logger, eventScheduler, horizon, frequency, payload, currentTimeFunc)

		// When
		err := metrolinkDeparturesDataLoadScheduler.Schedule(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given a valid configuration
When Schedule is called
And the EventScheduler dependency returns an error
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		horizon := time.Second * 30

		frequency := time.Second * 8

		payload := "foo"

		currentTimeFunc := time.Now

		eventScheduler := mock_repository.NewMockEventScheduler(ctrl)
		eventSchedulerErr := errors.New("FUBAR")
		eventScheduler.EXPECT().Schedule(ctx, gomock.Any()).Return(eventSchedulerErr)

		metrolinkDeparturesDataLoadScheduler, _ := v1.NewMetrolinkDeparturesDataLoadScheduler(logger, eventScheduler, horizon, frequency, payload, currentTimeFunc)

		// When
		err := metrolinkDeparturesDataLoadScheduler.Schedule(ctx)

		// Then
		assert.Equal(t, err, eventSchedulerErr)
	})
}
