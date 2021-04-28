package v1

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/Marchie/tf-experiment/lambda/internal/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type MetrolinkDeparturesDataLoadScheduler struct {
	logger          *zap.Logger
	eventScheduler  repository.EventScheduler
	horizon         time.Duration
	frequency       time.Duration
	payload         string
	currentTimeFunc func() time.Time
}

func NewMetrolinkDeparturesDataLoadScheduler(logger *zap.Logger, eventScheduler repository.EventScheduler, horizon time.Duration, frequency time.Duration, payload string, currentTimeFunc func() time.Time) (*MetrolinkDeparturesDataLoadScheduler, error) {
	if frequency == time.Duration(0) {
		return nil, errors.New("frequency cannot be 0")
	}

	return &MetrolinkDeparturesDataLoadScheduler{
		logger:          logger,
		eventScheduler:  eventScheduler,
		horizon:         horizon,
		frequency:       frequency,
		payload:         payload,
		currentTimeFunc: currentTimeFunc,
	}, nil
}

func (m *MetrolinkDeparturesDataLoadScheduler) Schedule(ctx context.Context) error {
	events := m.createEvents()

	return m.eventScheduler.Schedule(ctx, events)
}

func (m *MetrolinkDeparturesDataLoadScheduler) createEvents() []*domain.Event {
	now := m.currentTimeFunc()

	var events []*domain.Event

	for startInDuration := time.Duration(0); startInDuration < m.horizon; startInDuration += m.frequency {
		events = append(events, &domain.Event{
			StartTime: now.Add(startInDuration),
			Payload:   m.payload,
		})
	}

	return events
}
