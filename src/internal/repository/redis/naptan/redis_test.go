package naptan_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/redis/naptan"
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

func givenAtcoCodesFor940GZZMASTP(t *testing.T) []string {
	t.Helper()

	return []string{
		"9400ZZMASTP1",
		"9400ZZMASTP2",
		"9400ZZMASTP3",
		"9400ZZMASTP4",
	}
}

func givenAtcoCodesFor940GZZMAVIC(t *testing.T) []string {
	t.Helper()

	return []string{
		"9400ZZMAVIC1",
		"9400ZZMAVIC2",
		"9400ZZMAVIC3",
		"9400ZZMAVIC4",
	}
}

func givenAtcoCodesToStore(t *testing.T) map[string][]string {
	t.Helper()

	return map[string][]string{
		"940GZZMASTP": {
			"9400ZZMASTP1",
			"9400ZZMASTP2",
			"9400ZZMASTP3",
			"9400ZZMASTP4",
		},
		"940GZZMAVIC": {
			"9400ZZMAVIC1",
			"9400ZZMAVIC2",
			"9400ZZMAVIC3",
			"9400ZZMAVIC4",
		},
	}
}

func TestNaptanRedis_GetStopsInArea(t *testing.T) {
	t.Run(`Given a populated Redis stops in area repository
When GetStopsInArea is called with a StopAreaCode
Then AtcoCodes for that StopAreaCode are returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		stopAreaCode := "940GZZMASTP"

		expStopsInArea := givenAtcoCodesFor940GZZMASTP(t)

		var stopsInAreaFromRedis bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFromRedis).Encode(expStopsInArea); err != nil {
			t.Fatal(err)
		}

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "stopsinarea", stopAreaCode)).Return(stopsInAreaFromRedis.Bytes(), nil),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		stopsInArea, err := naptanRepository.GetStopsInArea(ctx, stopAreaCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, expStopsInArea, stopsInArea)
	})

	t.Run(`Given an error occurs getting a Redis connection from the connection pool
When GetStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreaCode := "940GZZMASTP"

		pool := mock_redis.NewMockPooler(ctrl)
		poolErr := errors.New("FUBAR")
		pool.EXPECT().GetContext(ctx).Return(nil, poolErr)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		stopsInArea, err := naptanRepository.GetStopsInArea(ctx, stopAreaCode)

		// Then
		assert.Nil(t, stopsInArea)
		assert.NotNil(t, err)
		assert.Equal(t, poolErr, err)
	})

	t.Run(`Given an error occurs getting data from the Redis repository
When GetStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreaCode := "940GZZMASTP"

		conn := mock_redis.NewMockConn(ctrl)
		connErr := errors.New("FUBAR")

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "stopsinarea", stopAreaCode)).Return(nil, connErr),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		stopsInArea, err := naptanRepository.GetStopsInArea(ctx, stopAreaCode)

		// Then
		assert.Nil(t, stopsInArea)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting stops in area for 940GZZMASTP: FUBAR")
	})

	t.Run(`Given an error occurs returning the Redis connection to the pool
When GetStopsInArea is called
Then an error is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreaCode := "940GZZMASTP"

		conn := mock_redis.NewMockConn(ctrl)
		connCloseErr := errors.New("FUBAR")

		expStopsInArea := givenAtcoCodesFor940GZZMASTP(t)

		var stopsInAreaFromRedis bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFromRedis).Encode(expStopsInArea); err != nil {
			t.Fatal(err)
		}

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "stopsinarea", stopAreaCode)).Return(stopsInAreaFromRedis.Bytes(), nil),
			conn.EXPECT().Close().Return(connCloseErr),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		stopsInArea, err := naptanRepository.GetStopsInArea(ctx, stopAreaCode)

		// Then
		assert.NotNil(t, stopsInArea)
		assert.Nil(t, err)

		assert.Equal(t, 1, observedLogs.Len())

		loggedItems := observedLogs.TakeAll()
		assert.Equal(t, zapcore.ErrorLevel, loggedItems[0].Level)
		assert.Equal(t, "error returning Redis connection to pool", loggedItems[0].Message)
		assert.Equal(t, "error", loggedItems[0].Context[0].Key)
		assert.Equal(t, connCloseErr, loggedItems[0].Context[0].Interface)
	})

	t.Run(`Given Redis returns invalid data
When GetStopsInArea is called
Then an error is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreaCode := "940GZZMASTP"

		conn := mock_redis.NewMockConn(ctrl)

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "stopsinarea", stopAreaCode)).Return([]byte("x"), nil),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		stopsInArea, err := naptanRepository.GetStopsInArea(ctx, stopAreaCode)

		// Then
		assert.Nil(t, stopsInArea)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error unmarshalling data for 940GZZMASTP: invalid character 'x' looking for beginning of value")
	})
}

func TestMetrolinkDeparturesRepository_Store(t *testing.T) {
	t.Run(`Given a map of StopAreaCodes to AtcoCodes
When StoreStopsInArea is called
Then the AtcoCodes are stored in Redis grouped by StopAreaCode`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreasToStore := givenAtcoCodesToStore(t)

		conn := mock_redis.NewMockConn(ctrl)

		var stopsInAreaFor940GZZMASTP bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMASTP).Encode(givenAtcoCodesFor940GZZMASTP(t)); err != nil {
			t.Fatal(err)
		}

		var stopsInAreaFor940GZZMAVIC bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMAVIC).Encode(givenAtcoCodesFor940GZZMAVIC(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "stopsinarea_940GZZMASTP", stopsInAreaFor940GZZMASTP.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "stopsinarea_940GZZMAVIC", stopsInAreaFor940GZZMAVIC.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil).Times(2)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		err := naptanRepository.StoreStopsInArea(ctx, stopAreasToStore)

		// Thens
		assert.Nil(t, err)
	})

	t.Run(`Given an error occurs getting a Redis connection from the connection pool
When StoreStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreasToStore := givenAtcoCodesToStore(t)

		pool := mock_redis.NewMockPooler(ctrl)
		poolErr := errors.New("FUBAR")
		pool.EXPECT().GetContext(ctx).Return(nil, poolErr)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		err := naptanRepository.StoreStopsInArea(ctx, stopAreasToStore)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, poolErr, err)
	})

	t.Run(`Given an error occurs sending data to Redis
When StoreStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreasToStore := givenAtcoCodesToStore(t)

		conn := mock_redis.NewMockConn(ctrl)

		var stopsInAreaFor940GZZMASTP bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMASTP).Encode(givenAtcoCodesFor940GZZMASTP(t)); err != nil {
			t.Fatal(err)
		}

		var stopsInAreaFor940GZZMAVIC bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMAVIC).Encode(givenAtcoCodesFor940GZZMAVIC(t)); err != nil {
			t.Fatal(err)
		}

		connErr := errors.New("FUBAR")
		conn.EXPECT().Send("SET", "stopsinarea_940GZZMASTP", stopsInAreaFor940GZZMASTP.String(), "PX", int64(15000)).Return(connErr)
		conn.EXPECT().Send("SET", "stopsinarea_940GZZMAVIC", stopsInAreaFor940GZZMAVIC.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil).Times(1)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		err := naptanRepository.StoreStopsInArea(ctx, stopAreasToStore)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error sending/flushing Redis connection: 1 error occurred:\n\t* error sending Redis command for StopAreaCode 940GZZMASTP: FUBAR\n\n")
	})

	t.Run(`Given an error occurs flushing the Redis connection
When StoreStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreasToStore := givenAtcoCodesToStore(t)

		conn := mock_redis.NewMockConn(ctrl)

		var stopsInAreaFor940GZZMASTP bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMASTP).Encode(givenAtcoCodesFor940GZZMASTP(t)); err != nil {
			t.Fatal(err)
		}

		var stopsInAreaFor940GZZMAVIC bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMAVIC).Encode(givenAtcoCodesFor940GZZMAVIC(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "stopsinarea_940GZZMASTP", stopsInAreaFor940GZZMASTP.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "stopsinarea_940GZZMAVIC", stopsInAreaFor940GZZMAVIC.String(), "PX", int64(15000)).Return(nil)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Flush().Return(connErr)
		conn.EXPECT().Receive().Return("OK", nil).Times(2)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		err := naptanRepository.StoreStopsInArea(ctx, stopAreasToStore)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error sending/flushing Redis connection: 1 error occurred:\n\t* error flushing Redis connection: FUBAR\n\n")
	})

	t.Run(`Given an error occurs receiving the response from the Redis connection
When StoreStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreasToStore := givenAtcoCodesToStore(t)

		conn := mock_redis.NewMockConn(ctrl)

		var stopsInAreaFor940GZZMASTP bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMASTP).Encode(givenAtcoCodesFor940GZZMASTP(t)); err != nil {
			t.Fatal(err)
		}

		var stopsInAreaFor940GZZMAVIC bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMAVIC).Encode(givenAtcoCodesFor940GZZMAVIC(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "stopsinarea_940GZZMASTP", stopsInAreaFor940GZZMASTP.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "stopsinarea_940GZZMAVIC", stopsInAreaFor940GZZMAVIC.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Receive().Return(nil, connErr)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		err := naptanRepository.StoreStopsInArea(ctx, stopAreasToStore)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error receiving on Redis connection: 1 error occurred:\n\t* FUBAR\n\n")
	})

	t.Run(`Given an error occurs returning the Redis connection to the pool
When StoreStopsInArea is called
Then an error message is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		keyPrefix := "stopsinarea"

		timeToLive := time.Second * 15

		stopAreasToStore := givenAtcoCodesToStore(t)

		conn := mock_redis.NewMockConn(ctrl)

		var stopsInAreaFor940GZZMASTP bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMASTP).Encode(givenAtcoCodesFor940GZZMASTP(t)); err != nil {
			t.Fatal(err)
		}

		var stopsInAreaFor940GZZMAVIC bytes.Buffer
		if err := json.NewEncoder(&stopsInAreaFor940GZZMAVIC).Encode(givenAtcoCodesFor940GZZMAVIC(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "stopsinarea_940GZZMASTP", stopsInAreaFor940GZZMASTP.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "stopsinarea_940GZZMAVIC", stopsInAreaFor940GZZMAVIC.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil).Times(2)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Close().Return(connErr)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		naptanRepository := naptan.NewNaptanRedis(logger, pool, keyPrefix, timeToLive)

		// When
		err := naptanRepository.StoreStopsInArea(ctx, stopAreasToStore)

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
