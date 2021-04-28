package cloudwatch

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core"
	"go.uber.org/zap"
)

type NaptanDataLoader struct {
	logger            *zap.Logger
	stopsInAreaLoader core.NaptanStopsInAreaLoader
}

func NewNaptanDataLoader(logger *zap.Logger, stopsInAreaLoader core.NaptanStopsInAreaLoader) *NaptanDataLoader {
	return &NaptanDataLoader{
		logger:            logger,
		stopsInAreaLoader: stopsInAreaLoader,
	}
}

func (n *NaptanDataLoader) Handler(ctx context.Context) error {
	return n.stopsInAreaLoader.LoadStopsInArea(ctx)
}
