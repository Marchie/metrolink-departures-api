package loader_test

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core/naptan/loader"
	mock_repository "github.com/Marchie/tf-experiment/lambda/internal/mocks/repository"
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

func mockStopsInAreaMap(t *testing.T) map[string][]string {
	t.Helper()

	stopsInAreaMap := make(map[string][]string)

	stopsInAreaMap["940GZZMASTP"] = []string{
		"9400ZZMASTP1",
		"9400ZZMASTP2",
		"9400ZZMASTP3",
		"9400ZZMASTP4",
	}

	return stopsInAreaMap
}

func TestStopsInAreaLoader_LoadStopsInArea(t *testing.T) {
	t.Run(`Given stops in area data can be fetched
When LoadStopsInArea is called
Then stops in area data is stored in the repository`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaMap := mockStopsInAreaMap(t)

		stopsInAreaFetcher := mock_repository.NewMockStopsInAreaFetcher(ctrl)
		stopsInAreaFetcher.EXPECT().FetchStopsInArea(ctx).Return(stopsInAreaMap, nil)

		stopsInAreaStorer := mock_repository.NewMockStopsInAreaStorer(ctrl)
		stopsInAreaStorer.EXPECT().StoreStopsInArea(ctx, stopsInAreaMap).Return(nil)

		stopsInAreaLoader := loader.NewStopsInAreaLoader(logger, stopsInAreaFetcher, stopsInAreaStorer)

		// When
		err := stopsInAreaLoader.LoadStopsInArea(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given stops in area data fails to fetch
When LoadStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaFetcher := mock_repository.NewMockStopsInAreaFetcher(ctrl)
		stopsInAreaFetcherErr := errors.New("FUBAR")
		stopsInAreaFetcher.EXPECT().FetchStopsInArea(ctx).Return(nil, stopsInAreaFetcherErr)

		stopsInAreaStorer := mock_repository.NewMockStopsInAreaStorer(ctrl)

		stopsInAreaLoader := loader.NewStopsInAreaLoader(logger, stopsInAreaFetcher, stopsInAreaStorer)

		// When
		err := stopsInAreaLoader.LoadStopsInArea(ctx)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, stopsInAreaFetcherErr, err)
	})

	t.Run(`Given stops in area data can be fetched
And an error occurs storing the data in the repository
When LoadStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaMap := mockStopsInAreaMap(t)

		stopsInAreaFetcher := mock_repository.NewMockStopsInAreaFetcher(ctrl)
		stopsInAreaFetcher.EXPECT().FetchStopsInArea(ctx).Return(stopsInAreaMap, nil)

		stopsInAreaStorer := mock_repository.NewMockStopsInAreaStorer(ctrl)
		stopsInAreaStorerErr := errors.New("FUBAR")
		stopsInAreaStorer.EXPECT().StoreStopsInArea(ctx, stopsInAreaMap).Return(stopsInAreaStorerErr)

		stopsInAreaLoader := loader.NewStopsInAreaLoader(logger, stopsInAreaFetcher, stopsInAreaStorer)

		// When
		err := stopsInAreaLoader.LoadStopsInArea(ctx)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, stopsInAreaStorerErr, err)
	})
}
