package cloudwatch

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core"
	"go.uber.org/zap"
)

type MetrolinkDeparturesDataLoadScheduler struct {
	logger         *zap.Logger
	eventScheduler core.EventScheduler
}

func NewMetrolinkDeparturesDataLoadScheduler(logger *zap.Logger, eventScheduler core.EventScheduler) *MetrolinkDeparturesDataLoadScheduler {
	return &MetrolinkDeparturesDataLoadScheduler{
		logger:         logger,
		eventScheduler: eventScheduler,
	}
}

func (m *MetrolinkDeparturesDataLoadScheduler) Handler(ctx context.Context) error {
	return m.eventScheduler.Schedule(ctx)
}
