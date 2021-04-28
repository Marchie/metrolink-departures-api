package loader

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/repository"
	"go.uber.org/zap"
)

type StopsInAreaLoader struct {
	logger             *zap.Logger
	stopsInAreaFetcher repository.StopsInAreaFetcher
	stopsInAreaStorer  repository.StopsInAreaStorer
}

func NewStopsInAreaLoader(logger *zap.Logger, stopsInAreaFetcher repository.StopsInAreaFetcher, stopsInAreaStorer repository.StopsInAreaStorer) *StopsInAreaLoader {
	return &StopsInAreaLoader{
		logger:             logger,
		stopsInAreaFetcher: stopsInAreaFetcher,
		stopsInAreaStorer:  stopsInAreaStorer,
	}
}

func (s *StopsInAreaLoader) LoadStopsInArea(ctx context.Context) error {
	stopsInAreaMap, err := s.stopsInAreaFetcher.FetchStopsInArea(ctx)
	if err != nil {
		return err
	}

	return s.stopsInAreaStorer.StoreStopsInArea(ctx, stopsInAreaMap)
}
