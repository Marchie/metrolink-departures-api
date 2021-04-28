package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

type Pooler interface {
	Get() redis.Conn
	GetContext(ctx context.Context) (redis.Conn, error)
	Stats() redis.PoolStats
	ActiveCount() int
	IdleCount() int
	Close() error
}
