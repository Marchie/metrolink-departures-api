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

func TestMetrolinkDeparturesDataLoadScheduler_Handler(t *testing.T) {
	t.Run(`Given valid configuration
When Handler is called
Then data loading is scheduled`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		eventScheduler := mock_core.NewMockEventScheduler(ctrl)
		eventScheduler.EXPECT().Schedule(ctx).Return(nil)

		metrolinkDeparturesDataLoadScheduler := cloudwatch.NewMetrolinkDeparturesDataLoadScheduler(logger, eventScheduler)

		// When
		err := metrolinkDeparturesDataLoadScheduler.Handler(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given an event scheduler error
When Handler is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		eventScheduler := mock_core.NewMockEventScheduler(ctrl)
		eventSchedulerErr := errors.New("FUBAR")
		eventScheduler.EXPECT().Schedule(ctx).Return(eventSchedulerErr)

		metrolinkDeparturesDataLoadScheduler := cloudwatch.NewMetrolinkDeparturesDataLoadScheduler(logger, eventScheduler)

		// When
		err := metrolinkDeparturesDataLoadScheduler.Handler(ctx)

		// Then
		assert.Equal(t, err, eventSchedulerErr)
	})
}
