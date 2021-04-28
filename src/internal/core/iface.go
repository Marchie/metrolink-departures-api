package core

import (
	"context"
	"io"
)

type StopAreaDeparturesJsoner interface {
	Json(ctx context.Context, stopAreaCode string) (io.ReadCloser, int, error)
}

type EventScheduler interface {
	Schedule(ctx context.Context) error
}

type MetrolinkDeparturesLoader interface {
	Load(ctx context.Context) error
}

type NaptanStopsInAreaLoader interface {
	LoadStopsInArea(ctx context.Context) error
}
