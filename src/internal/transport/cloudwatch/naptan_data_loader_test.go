package cloudwatch_test

import (
	"context"
	mock_core "github.com/Marchie/tf-experiment/lambda/internal/mocks/core"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/cloudwatch"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNaptanDataLoader_Handler(t *testing.T) {
	t.Run(`Given valid configuration
When Handler is called
Then NaPTAN data is loaded`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		naptanStopsInAreaLoader := mock_core.NewMockNaptanStopsInAreaLoader(ctrl)
		naptanStopsInAreaLoader.EXPECT().LoadStopsInArea(ctx).Return(nil)

		naptanDataLoader := cloudwatch.NewNaptanDataLoader(logger, naptanStopsInAreaLoader)

		// When
		err := naptanDataLoader.Handler(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given an error occurs loading NaPTAN data
When Handler is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		naptanStopsInAreaLoader := mock_core.NewMockNaptanStopsInAreaLoader(ctrl)
		naptanStopsInAreaLoaderErr := errors.New("FUBAR")
		naptanStopsInAreaLoader.EXPECT().LoadStopsInArea(ctx).Return(naptanStopsInAreaLoaderErr)

		naptanDataLoader := cloudwatch.NewNaptanDataLoader(logger, naptanStopsInAreaLoader)

		// When
		err := naptanDataLoader.Handler(ctx)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, naptanStopsInAreaLoaderErr, err)
	})
}
