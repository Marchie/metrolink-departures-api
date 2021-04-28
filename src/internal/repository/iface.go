package repository

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"time"
)

type AtcoCodeLister interface {
	GetAtcoCodesInStopArea(stopAreaCode string) ([]string, error)
}

type EventScheduler interface {
	Schedule(ctx context.Context, events []*domain.Event) error
}

type MetrolinkDeparturesFetcher interface {
	Fetch(ctx context.Context) (*domain.MetrolinkDepartures, error)
}

type MetrolinkDeparturesGetter interface {
	Get(ctx context.Context, stopAreaCode string) ([]*domain.MetrolinkDeparture, error)
}

type MetrolinkDeparturesStorer interface {
	Store(ctx context.Context, departures []*domain.MetrolinkDeparture) error
}

type SystemStatusGetter interface {
	Get(ctx context.Context) (*time.Time, error)
}

type SystemStatusSetter interface {
	Set(ctx context.Context, lastUpdated time.Time) error
}

type PlatformNamer interface {
	GetPlatformNameForAtcoCode(atcoCode string) (*string, error)
}

type StopsInAreaFetcher interface {
	FetchStopsInArea(ctx context.Context) (map[string][]string, error)
}

type StopsInAreaGetter interface {
	GetStopsInArea(ctx context.Context, stopAreaCode string) ([]string, error)
}

type StopsInAreaStorer interface {
	StoreStopsInArea(ctx context.Context, stopsInArea map[string][]string) error
}
