package naptan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	redis2 "github.com/Marchie/tf-experiment/lambda/pkg/redis"
	"github.com/gomodule/redigo/redis"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type NaptanRedis struct {
	logger     *zap.Logger
	pool       redis2.Pooler
	keyPrefix  string
	timeToLive time.Duration
}

func NewNaptanRedis(logger *zap.Logger, pool redis2.Pooler, keyPrefix string, timeToLive time.Duration) *NaptanRedis {
	return &NaptanRedis{
		logger:     logger,
		pool:       pool,
		keyPrefix:  keyPrefix,
		timeToLive: timeToLive,
	}
}

func (n *NaptanRedis) GetStopsInArea(ctx context.Context, stopAreaCode string) ([]string, error) {
	conn, err := n.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			n.logger.Error("error returning Redis connection to pool", zap.Error(err))
		}
	}()

	stopsInAreaJson, err := redis.Bytes(conn.Do("GET", n.key(stopAreaCode)))
	if err != nil {
		return nil, errors.Wrapf(err, "error getting stops in area for %s", stopAreaCode)
	}

	var atcoCodes []string
	if err := json.Unmarshal(stopsInAreaJson, &atcoCodes); err != nil {
		return nil, errors.Wrapf(err, "error unmarshalling data for %s", stopAreaCode)
	}

	return atcoCodes, nil
}

func (n *NaptanRedis) StoreStopsInArea(ctx context.Context, stopsInArea map[string][]string) error {
	conn, err := n.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			n.logger.Error("error returning Redis connection to pool", zap.Error(err))
		}
	}()

	chReceive, chErr := n.send(conn, stopsInArea)
	if err := n.receive(conn, chReceive); err != nil {
		return errors.Wrap(err, "error receiving on Redis connection")
	}

	if err := <-chErr; err != nil {
		return errors.Wrap(err, "error sending/flushing Redis connection")
	}

	return nil
}

func (n *NaptanRedis) send(conn redis.Conn, stopsInArea map[string][]string) (chan int, chan error) {
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

		for stopAreaCode, atcoCodes := range stopsInArea {
			var stopsInAreaJson bytes.Buffer
			if err := json.NewEncoder(&stopsInAreaJson).Encode(atcoCodes); err != nil {
				errs = multierror.Append(errs, errors.Wrapf(err, "error encoding AtcoCodes for StopAreaCode %s", stopAreaCode))
				continue
			}

			if err := conn.Send("SET", n.key(stopAreaCode), stopsInAreaJson.String(), "PX", n.timeToLive.Milliseconds()); err != nil {
				errs = multierror.Append(errs, errors.Wrapf(err, "error sending Redis command for StopAreaCode %s", stopAreaCode))
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

func (n *NaptanRedis) receive(conn redis.Conn, chReceive <-chan int) error {
	var errs error

	limit := <-chReceive

	for i := 0; i < limit; i++ {
		if _, err := conn.Receive(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (n NaptanRedis) key(stopAreaCode string) string {
	return fmt.Sprintf("%s_%s", n.keyPrefix, stopAreaCode)
}
