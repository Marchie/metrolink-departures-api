package v1

import (
	"context"
	redis2 "github.com/Marchie/tf-experiment/lambda/pkg/redis"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

type MetrolinkDeparturesSystemStatusRepository struct {
	logger          *zap.Logger
	pool            redis2.Pooler
	systemStatusKey string
}

func NewMetrolinkDeparturesSystemStatusRepository(logger *zap.Logger, pool redis2.Pooler, systemStatusKey string) *MetrolinkDeparturesSystemStatusRepository {
	return &MetrolinkDeparturesSystemStatusRepository{
		logger:          logger,
		pool:            pool,
		systemStatusKey: systemStatusKey,
	}
}

func (m *MetrolinkDeparturesSystemStatusRepository) Get(ctx context.Context) (*time.Time, error) {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			m.logger.Error("error returning Redis connection to pool", zap.Error(err))
		}
	}()

	lastUpdatedStr, err := redis.String(conn.Do("GET", m.systemStatusKey))
	if err != nil {
		return nil, err
	}

	lastUpdated, err := time.Parse(time.RFC3339, lastUpdatedStr)
	if err != nil {
		return nil, err
	}

	return &lastUpdated, nil
}

func (m *MetrolinkDeparturesSystemStatusRepository) Set(ctx context.Context, lastUpdated time.Time) error {
	conn, err := m.pool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			m.logger.Error("error returning Redis connection to pool", zap.Error(err))
		}
	}()

	_, err = conn.Do("SET", m.systemStatusKey, lastUpdated.Format(time.RFC3339))
	return err
}
