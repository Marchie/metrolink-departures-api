package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	redis2 "github.com/Marchie/tf-experiment/lambda/pkg/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type MetrolinkDeparturesRepository struct {
	logger               *zap.Logger
	pool                 redis2.Pooler
	departuresKeyPrefix  string
	departuresTimeToLive time.Duration
}

func NewMetrolinkDeparturesRepository(logger *zap.Logger, pool redis2.Pooler, departuresKeyPrefix string, departuresTimeToLive time.Duration) *MetrolinkDeparturesRepository {
	return &MetrolinkDeparturesRepository{
		logger:               logger,
		pool:                 pool,
		departuresKeyPrefix:  departuresKeyPrefix,
		departuresTimeToLive: departuresTimeToLive,
	}
}

func (m *MetrolinkDeparturesRepository) Get(ctx context.Context, atcoCode string) ([]*domain.MetrolinkDeparture, error) {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			m.logger.Error("error returning Redis connection to pool", zap.Error(err))
		}
	}()

	departuresJson, err := redis.Bytes(conn.Do("GET", m.departuresKey(atcoCode)))
	if err != nil {
		return nil, err
	}

	var departures []*domain.MetrolinkDeparture

	if err := json.Unmarshal(departuresJson, &departures); err != nil {
		return nil, err
	}

	return departures, nil
}

func (m *MetrolinkDeparturesRepository) Store(ctx context.Context, departures []*domain.MetrolinkDeparture) error {
	groupedDeparturesByAtcoCode := m.groupDeparturesByAtcoCode(departures)

	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			m.logger.Error("error returning Redis connection to pool", zap.Error(err))
		}
	}()

	chReceive, chErr := m.send(conn, groupedDeparturesByAtcoCode)
	if err := m.receive(conn, chReceive); err != nil {
		return errors.Wrap(err, "error receiving on Redis connection")
	}

	if err := <-chErr; err != nil {
		return errors.Wrap(err, "error sending/flushing Redis connection")
	}

	return nil
}

func (*MetrolinkDeparturesRepository) groupDeparturesByAtcoCode(departures []*domain.MetrolinkDeparture) map[string][]*domain.MetrolinkDeparture {
	groupedDepartures := make(map[string][]*domain.MetrolinkDeparture)

	for _, departure := range departures {
		groupedDepartures[departure.AtcoCode] = append(groupedDepartures[departure.AtcoCode], departure)
	}

	return groupedDepartures
}

func (m *MetrolinkDeparturesRepository) send(conn redis.Conn, groupedDeparturesByAtcoCode map[string][]*domain.MetrolinkDeparture) (chan int, chan error) {
	chReceive := make(chan int)
	chErr := make(chan error, 1)

	go func() {
		defer close(chErr)
		defer close(chReceive)

		var errs error

		defer func() {
			chErr <- errs
		}()

		i := 0

		for atcoCode, departures := range groupedDeparturesByAtcoCode {
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(departures); err != nil {
				errs = multierror.Append(errs, errors.Wrapf(err, "error encoding departures as JSON for AtcoCode %s", atcoCode))
				continue
			}

			if err := conn.Send("SET", m.departuresKey(atcoCode), buf.String(), "PX", m.departuresTimeToLive.Milliseconds()); err != nil {
				errs = multierror.Append(errs, errors.Wrapf(err, "error sending Redis command for AtcoCode %s", atcoCode))
				continue
			}

			i++
		}

		if err := conn.Flush(); err != nil {
			errs = multierror.Append(errs, errors.Wrap(err, "error flushing Redis connection"))
		}

		chReceive <- i
	}()

	return chReceive, chErr
}

func (*MetrolinkDeparturesRepository) receive(conn redis.Conn, chReceive <-chan int) error {
	var errs error

	limit := <-chReceive

	for i := 0; i < limit; i++ {
		if _, err := conn.Receive(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (m *MetrolinkDeparturesRepository) departuresKey(atcoCode string) string {
	return fmt.Sprintf("%s_%s", m.departuresKeyPrefix, atcoCode)
}
