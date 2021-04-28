package v1_test

import (
	"context"
	v1 "github.com/Marchie/tf-experiment/lambda/internal/repository/redis/system/status/v1"
	mock_redis "github.com/Marchie/tf-experiment/lambda/pkg/mocks/redis"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"testing"
	"time"
)

func mockLogger(t *testing.T) *zap.Logger {
	t.Helper()

	zapCore, _ := observer.New(zapcore.DebugLevel)
	return zap.New(zapCore)
}

func TestMetrolinkDeparturesSystemStatusRepository_Get(t *testing.T) {
	t.Run(`Given a populated Redis status in the repository
When Get is called
Then the last updated time is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		expLastUpdatedTime := time.Date(2021, time.April, 24, 2, 30, 15, 0, time.UTC)

		conn := mock_redis.NewMockConn(ctrl)

		gomock.InOrder(
			conn.EXPECT().Do("GET", systemStatusKey).Return(expLastUpdatedTime.Format(time.RFC3339), nil),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		lastUpdatedTime, err := metrolinkDeparturesRepository.Get(ctx)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, &expLastUpdatedTime, lastUpdatedTime)
	})

	t.Run(`Given an error occurs getting a Redis connection from the connection pool
When Get is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		pool := mock_redis.NewMockPooler(ctrl)
		poolErr := errors.New("FUBAR")
		pool.EXPECT().GetContext(ctx).Return(nil, poolErr)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		lastUpdatedTime, err := metrolinkDeparturesRepository.Get(ctx)

		// Then
		assert.Nil(t, lastUpdatedTime)
		assert.NotNil(t, err)
		assert.Equal(t, poolErr, err)
	})

	t.Run(`Given an error occurs getting data from the Redis repository
When Get is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		conn := mock_redis.NewMockConn(ctrl)
		connErr := errors.New("FUBAR")

		gomock.InOrder(
			conn.EXPECT().Do("GET", systemStatusKey).Return(nil, connErr),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		lastUpdatedTime, err := metrolinkDeparturesRepository.Get(ctx)

		// Then
		assert.Nil(t, lastUpdatedTime)
		assert.NotNil(t, err)
		assert.Equal(t, connErr, err)
	})

	t.Run(`Given an error occurs returning the Redis connection to the pool
When Get is called
Then an error is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		systemStatusKey := "status"

		conn := mock_redis.NewMockConn(ctrl)
		connCloseErr := errors.New("FUBAR")

		expLastUpdatedTime := time.Date(2021, time.April, 24, 2, 30, 15, 0, time.UTC)

		gomock.InOrder(
			conn.EXPECT().Do("GET", systemStatusKey).Return(expLastUpdatedTime.Format(time.RFC3339), nil),
			conn.EXPECT().Close().Return(connCloseErr),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		lastUpdatedTime, err := metrolinkDeparturesRepository.Get(ctx)

		// Then
		assert.NotNil(t, lastUpdatedTime)
		assert.Nil(t, err)

		assert.Equal(t, 1, observedLogs.Len())

		loggedItems := observedLogs.TakeAll()
		assert.Equal(t, zapcore.ErrorLevel, loggedItems[0].Level)
		assert.Equal(t, "error returning Redis connection to pool", loggedItems[0].Message)
		assert.Equal(t, "error", loggedItems[0].Context[0].Key)
		assert.Equal(t, connCloseErr, loggedItems[0].Context[0].Interface)
	})

	t.Run(`Given Redis returns invalid data
When Get is called
Then an error is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		conn := mock_redis.NewMockConn(ctrl)

		gomock.InOrder(
			conn.EXPECT().Do("GET", systemStatusKey).Return("FUBAR", nil),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		lastUpdatedTime, err := metrolinkDeparturesRepository.Get(ctx)

		// Then
		assert.Nil(t, lastUpdatedTime)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "parsing time \"FUBAR\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"FUBAR\" as \"2006\"")
	})
}

func TestMetrolinkDeparturesRepository_SetStatus(t *testing.T) {
	t.Run(`Given a last updated time value
When Set is called
Then the last updated time is stored in the repository`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		lastUpdatedTime := time.Date(2021, time.April, 24, 2, 24, 44, 0, time.UTC)

		conn := mock_redis.NewMockConn(ctrl)

		conn.EXPECT().Do("SET", systemStatusKey, lastUpdatedTime.Format(time.RFC3339)).Return(nil, nil)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		err := metrolinkDeparturesRepository.Set(ctx, lastUpdatedTime)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given an error occurs getting a Redis connection from the connection pool
When Set is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		lastUpdatedTime := time.Date(2021, time.April, 24, 2, 24, 44, 0, time.UTC)

		pool := mock_redis.NewMockPooler(ctrl)
		poolErr := errors.New("FUBAR")
		pool.EXPECT().GetContext(ctx).Return(nil, poolErr)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		err := metrolinkDeparturesRepository.Set(ctx, lastUpdatedTime)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, poolErr, err)
	})

	t.Run(`Given an error occurs sending data to Redis
When Set is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		systemStatusKey := "status"

		lastUpdatedTime := time.Date(2021, time.April, 24, 2, 24, 44, 0, time.UTC)

		conn := mock_redis.NewMockConn(ctrl)

		connErr := errors.New("FUBAR")
		conn.EXPECT().Do("SET", systemStatusKey, lastUpdatedTime.Format(time.RFC3339)).Return(nil, connErr)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		err := metrolinkDeparturesRepository.Set(ctx, lastUpdatedTime)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "FUBAR")
	})

	t.Run(`Given an error occurs returning the Redis connection to the pool
When Set is called
Then an error message is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		systemStatusKey := "status"

		lastUpdatedTime := time.Date(2021, time.April, 24, 2, 24, 44, 0, time.UTC)

		conn := mock_redis.NewMockConn(ctrl)

		conn.EXPECT().Do("SET", systemStatusKey, lastUpdatedTime.Format(time.RFC3339)).Return(nil, nil)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Close().Return(connErr)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesSystemStatusRepository(logger, pool, systemStatusKey)

		// When
		err := metrolinkDeparturesRepository.Set(ctx, lastUpdatedTime)

		// Then
		assert.Nil(t, err)

		assert.Equal(t, 1, observedLogs.Len())

		loggedItems := observedLogs.TakeAll()
		assert.Equal(t, zapcore.ErrorLevel, loggedItems[0].Level)
		assert.Equal(t, "error returning Redis connection to pool", loggedItems[0].Message)
		assert.Equal(t, "error", loggedItems[0].Context[0].Key)
		assert.Equal(t, connErr, loggedItems[0].Context[0].Interface)
	})
}
