package loader_test

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core/departures/metrolink/loader"
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

func givenCurrentTimeFunction(t *testing.T) func() time.Time {
	t.Helper()

	return func() time.Time {
		return time.Date(2021, time.March, 30, 22, 31, 18, 0, time.UTC)
	}
}

func givenStaleDataThreshold(t *testing.T) time.Duration {
	t.Helper()

	return time.Second * 45
}

func givenLastUpdatedTimeWithinThreshold(t *testing.T) time.Time {
	t.Helper()

	return time.Date(2021, time.March, 30, 22, 31, 15, 0, time.UTC)
}

func givenLastUpdatedTimeOutsideOfThreshold(t *testing.T) time.Time {
	t.Helper()

	return time.Date(2021, time.March, 30, 22, 30, 32, 0, time.UTC)
}

func givenMetrolinkDeparturesFromSource(t *testing.T) *domain.MetrolinkDepartures {
	t.Helper()

	return &domain.MetrolinkDepartures{
		Departures: []*domain.MetrolinkDeparture{
			{
				AtcoCode:    "9400ZZMASTP1",
				Order:       0,
				Destination: "Victoria",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "2",
				Platform:    nil,
				LastUpdated: givenLastUpdatedTimeWithinThreshold(t),
			},
			{
				AtcoCode:    "9400ZZMASTP1",
				Order:       1,
				Destination: "Rochdale",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "5",
				Platform:    nil,
				LastUpdated: givenLastUpdatedTimeWithinThreshold(t),
			},
		},
		LastUpdated: givenLastUpdatedTimeWithinThreshold(t),
	}
}

func givenStaleMetrolinkDeparturesFromSource(t *testing.T) *domain.MetrolinkDepartures {
	t.Helper()

	return &domain.MetrolinkDepartures{
		Departures: []*domain.MetrolinkDeparture{
			{
				AtcoCode:    "9400ZZMASTP1",
				Order:       0,
				Destination: "Victoria",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "2",
				Platform:    nil,
				LastUpdated: givenLastUpdatedTimeOutsideOfThreshold(t),
			},
			{
				AtcoCode:    "9400ZZMASTP1",
				Order:       1,
				Destination: "Rochdale",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "5",
				Platform:    nil,
				LastUpdated: givenLastUpdatedTimeOutsideOfThreshold(t),
			},
		},
		LastUpdated: givenLastUpdatedTimeOutsideOfThreshold(t),
	}
}

func givenMetrolinkDeparturesFromSourceWithPlatformsExpectation(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	platform := "D"

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       0,
			Destination: "Victoria",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "2",
			Platform:    &platform,
			LastUpdated: givenLastUpdatedTimeWithinThreshold(t),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       1,
			Destination: "Rochdale",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "5",
			Platform:    &platform,
			LastUpdated: givenLastUpdatedTimeWithinThreshold(t),
		},
	}
}

func TestMetrolinkDeparturesLoader_Load(t *testing.T) {
	t.Run(`Given Metrolink Departures from a source
When Load is executed
Then the departures are stored in a repository
And the system status is stored in a repository`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresFromSource := givenMetrolinkDeparturesFromSource(t)

		fetcher := mock_repository.NewMockMetrolinkDeparturesFetcher(ctrl)
		fetcher.EXPECT().Fetch(ctx).Return(departuresFromSource, nil)

		platform := "D"

		platformNamer := mock_repository.NewMockPlatformNamer(ctrl)
		platformNamer.EXPECT().GetPlatformNameForAtcoCode("9400ZZMASTP1").Times(2).Return(&platform, nil)

		departuresFromSourceWithPlatformsExpectation := givenMetrolinkDeparturesFromSourceWithPlatformsExpectation(t)
		departuresStorer := mock_repository.NewMockMetrolinkDeparturesStorer(ctrl)
		departuresStorer.EXPECT().Store(ctx, departuresFromSourceWithPlatformsExpectation).Return(nil)

		systemStatusSetter := mock_repository.NewMockSystemStatusSetter(ctrl)
		systemStatusSetter.EXPECT().Set(ctx, departuresFromSource.LastUpdated).Return(nil)

		currentTimeFunc := givenCurrentTimeFunction(t)
		staleDataThreshold := givenStaleDataThreshold(t)

		metrolinkDeparturesLoader := loader.NewMetrolinkDeparturesLoader(logger, fetcher, platformNamer, departuresStorer, systemStatusSetter, currentTimeFunc, staleDataThreshold)

		// When
		err := metrolinkDeparturesLoader.Load(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given Metrolink Departures cannot be fetched from a source
When Load is executed
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		fetcher := mock_repository.NewMockMetrolinkDeparturesFetcher(ctrl)
		fetcherErr := errors.New("FUBAR")
		fetcher.EXPECT().Fetch(ctx).Return(nil, fetcherErr)

		platformNamer := mock_repository.NewMockPlatformNamer(ctrl)

		departuresStorer := mock_repository.NewMockMetrolinkDeparturesStorer(ctrl)
		systemStatusSetter := mock_repository.NewMockSystemStatusSetter(ctrl)

		currentTimeFunc := givenCurrentTimeFunction(t)
		staleDataThreshold := givenStaleDataThreshold(t)

		metrolinkDeparturesLoader := loader.NewMetrolinkDeparturesLoader(logger, fetcher, platformNamer, departuresStorer, systemStatusSetter, currentTimeFunc, staleDataThreshold)

		// When
		err := metrolinkDeparturesLoader.Load(ctx)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, err, fetcherErr)
	})

	t.Run(`Given Metrolink departures are stale
When Load is executed
Then an error is logged
And the system status is stored
And stale records are not stored`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		departuresFromSource := givenStaleMetrolinkDeparturesFromSource(t)

		fetcher := mock_repository.NewMockMetrolinkDeparturesFetcher(ctrl)
		fetcher.EXPECT().Fetch(ctx).Return(departuresFromSource, nil)

		platformNamer := mock_repository.NewMockPlatformNamer(ctrl)

		departuresStorer := mock_repository.NewMockMetrolinkDeparturesStorer(ctrl)
		departuresStorer.EXPECT().Store(ctx, nil).Return(nil)

		systemStatusSetter := mock_repository.NewMockSystemStatusSetter(ctrl)
		systemStatusSetter.EXPECT().Set(ctx, departuresFromSource.LastUpdated).Return(nil)

		currentTimeFunc := givenCurrentTimeFunction(t)
		staleDataThreshold := givenStaleDataThreshold(t)

		metrolinkDeparturesLoader := loader.NewMetrolinkDeparturesLoader(logger, fetcher, platformNamer, departuresStorer, systemStatusSetter, currentTimeFunc, staleDataThreshold)

		// When
		err := metrolinkDeparturesLoader.Load(ctx)

		// Then
		assert.Nil(t, err)

		assert.Equal(t, 2, observedLogs.Len())
		loggedItems := observedLogs.TakeAll()
		for i := 0; i < 2; i++ {
			assert.Equal(t, zapcore.ErrorLevel, loggedItems[i].Level)
			assert.Equal(t, "error with source data - stale data received", loggedItems[i].Message)
			assert.Equal(t, "atcoCode", loggedItems[i].Context[0].Key)
			assert.Equal(t, "9400ZZMASTP1", loggedItems[i].Context[0].String)
			assert.Equal(t, "lastUpdated", loggedItems[i].Context[1].Key)
			assert.Equal(t, givenLastUpdatedTimeOutsideOfThreshold(t).UnixNano(), loggedItems[i].Context[1].Integer)
			assert.Equal(t, "staleDataThreshold", loggedItems[i].Context[2].Key)
			assert.Equal(t, staleDataThreshold.Nanoseconds(), loggedItems[i].Context[2].Integer)
			assert.Equal(t, "ageOfData", loggedItems[i].Context[3].Key)
			assert.Equal(t, (time.Second * 46).Nanoseconds(), loggedItems[i].Context[3].Integer)
		}
	})

	t.Run(`Given platform name cannot be queried for Metrolink Departures
When Load is executed
Then an error is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		departuresFromSource := givenMetrolinkDeparturesFromSource(t)

		fetcher := mock_repository.NewMockMetrolinkDeparturesFetcher(ctrl)
		fetcher.EXPECT().Fetch(ctx).Return(departuresFromSource, nil)

		platformNamer := mock_repository.NewMockPlatformNamer(ctrl)
		platformNamerErr := errors.New("FUBAR")
		platformNamer.EXPECT().GetPlatformNameForAtcoCode("9400ZZMASTP1").Times(2).Return(nil, platformNamerErr)

		departuresStorer := mock_repository.NewMockMetrolinkDeparturesStorer(ctrl)
		departuresStorer.EXPECT().Store(ctx, departuresFromSource.Departures).Return(nil)

		systemStatusSetter := mock_repository.NewMockSystemStatusSetter(ctrl)
		systemStatusSetter.EXPECT().Set(ctx, departuresFromSource.LastUpdated)

		currentTimeFunc := givenCurrentTimeFunction(t)
		staleDataThreshold := givenStaleDataThreshold(t)

		metrolinkDeparturesLoader := loader.NewMetrolinkDeparturesLoader(logger, fetcher, platformNamer, departuresStorer, systemStatusSetter, currentTimeFunc, staleDataThreshold)

		// When
		err := metrolinkDeparturesLoader.Load(ctx)

		// Then
		assert.Nil(t, err)

		assert.Equal(t, 2, observedLogs.Len())
		loggedItems := observedLogs.TakeAll()
		for i := 0; i < 2; i++ {
			assert.Equal(t, zapcore.ErrorLevel, loggedItems[i].Level)
			assert.Equal(t, "error getting platform name", loggedItems[i].Message)
			assert.Equal(t, "error", loggedItems[i].Context[0].Key)
			assert.Equal(t, platformNamerErr, loggedItems[i].Context[0].Interface)
			assert.Equal(t, "atcoCode", loggedItems[i].Context[1].Key)
			assert.Equal(t, "9400ZZMASTP1", loggedItems[i].Context[1].String)
		}
	})

	t.Run(`Given system status cannot be stored
When Load is executed
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresFromSource := givenMetrolinkDeparturesFromSource(t)

		fetcher := mock_repository.NewMockMetrolinkDeparturesFetcher(ctrl)
		fetcher.EXPECT().Fetch(ctx).Return(departuresFromSource, nil)

		platformNamer := mock_repository.NewMockPlatformNamer(ctrl)

		departuresStorer := mock_repository.NewMockMetrolinkDeparturesStorer(ctrl)

		systemStatusErr := errors.New("FUBAR")
		systemStatusSetter := mock_repository.NewMockSystemStatusSetter(ctrl)
		systemStatusSetter.EXPECT().Set(ctx, departuresFromSource.LastUpdated).Return(systemStatusErr)

		currentTimeFunc := givenCurrentTimeFunction(t)
		staleDataThreshold := givenStaleDataThreshold(t)

		metrolinkDeparturesLoader := loader.NewMetrolinkDeparturesLoader(logger, fetcher, platformNamer, departuresStorer, systemStatusSetter, currentTimeFunc, staleDataThreshold)

		// When
		err := metrolinkDeparturesLoader.Load(ctx)

		// Then
		assert.Equal(t, systemStatusErr, err)
	})

	t.Run(`Given Metrolink Departures cannot be stored
When Load is executed
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresFromSource := givenMetrolinkDeparturesFromSource(t)

		fetcher := mock_repository.NewMockMetrolinkDeparturesFetcher(ctrl)
		fetcher.EXPECT().Fetch(ctx).Return(departuresFromSource, nil)

		platform := "D"

		platformNamer := mock_repository.NewMockPlatformNamer(ctrl)
		platformNamer.EXPECT().GetPlatformNameForAtcoCode("9400ZZMASTP1").Times(2).Return(&platform, nil)

		departuresStorerErr := errors.New("FUBAR")
		departuresFromSourceWithPlatformsExpectation := givenMetrolinkDeparturesFromSourceWithPlatformsExpectation(t)
		departuresStorer := mock_repository.NewMockMetrolinkDeparturesStorer(ctrl)
		departuresStorer.EXPECT().Store(ctx, departuresFromSourceWithPlatformsExpectation).Return(departuresStorerErr)

		systemStatusSetter := mock_repository.NewMockSystemStatusSetter(ctrl)
		systemStatusSetter.EXPECT().Set(ctx, departuresFromSource.LastUpdated).Return(nil)

		currentTimeFunc := givenCurrentTimeFunction(t)
		staleDataThreshold := givenStaleDataThreshold(t)

		metrolinkDeparturesLoader := loader.NewMetrolinkDeparturesLoader(logger, fetcher, platformNamer, departuresStorer, systemStatusSetter, currentTimeFunc, staleDataThreshold)

		// When
		err := metrolinkDeparturesLoader.Load(ctx)

		// Then
		assert.Equal(t, departuresStorerErr, err)
	})
}
