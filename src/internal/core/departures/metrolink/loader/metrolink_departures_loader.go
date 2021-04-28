package loader

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/Marchie/tf-experiment/lambda/internal/repository"
	"go.uber.org/zap"
	"time"
)

type MetrolinkDeparturesLoader struct {
	logger             *zap.Logger
	departuresSource   repository.MetrolinkDeparturesFetcher
	platformNamer      repository.PlatformNamer
	departuresStorer   repository.MetrolinkDeparturesStorer
	systemStatusSetter repository.SystemStatusSetter
	currentTimeFunc    func() time.Time
	staleDataThreshold time.Duration
}

func NewMetrolinkDeparturesLoader(logger *zap.Logger, departuresSource repository.MetrolinkDeparturesFetcher, platformNamer repository.PlatformNamer, departuresStorer repository.MetrolinkDeparturesStorer, systemStatusSetter repository.SystemStatusSetter, currentTimeFunc func() time.Time, staleDataThreshold time.Duration) *MetrolinkDeparturesLoader {
	return &MetrolinkDeparturesLoader{
		logger:             logger,
		departuresSource:   departuresSource,
		platformNamer:      platformNamer,
		departuresStorer:   departuresStorer,
		systemStatusSetter: systemStatusSetter,
		currentTimeFunc:    currentTimeFunc,
		staleDataThreshold: staleDataThreshold,
	}
}

func (m *MetrolinkDeparturesLoader) Load(ctx context.Context) error {
	departuresFromSource, err := m.departuresSource.Fetch(ctx)
	if err != nil {
		return err
	}

	if err := m.systemStatusSetter.Set(ctx, departuresFromSource.LastUpdated); err != nil {
		return err
	}

	var departuresToStore []*domain.MetrolinkDeparture

	for _, departure := range departuresFromSource.Departures {
		if m.dataIsStale(departure) {
			m.logger.Error("error with source data - stale data received", zap.String("atcoCode", departure.AtcoCode), zap.Time("lastUpdated", departure.LastUpdated), zap.Duration("staleDataThreshold", m.staleDataThreshold), zap.Duration("ageOfData", m.currentTimeFunc().Sub(departure.LastUpdated)))
			continue
		}

		platform, err := m.platformNamer.GetPlatformNameForAtcoCode(departure.AtcoCode)
		if err != nil {
			m.logger.Error("error getting platform name", zap.Error(err), zap.String("atcoCode", departure.AtcoCode))
		}

		departure.Platform = platform

		departuresToStore = append(departuresToStore, departure)
	}

	return m.departuresStorer.Store(ctx, departuresToStore)
}

func (m *MetrolinkDeparturesLoader) dataIsStale(departure *domain.MetrolinkDeparture) bool {
	return departure.LastUpdated.Before(m.currentTimeFunc().Add(-m.staleDataThreshold))
}
