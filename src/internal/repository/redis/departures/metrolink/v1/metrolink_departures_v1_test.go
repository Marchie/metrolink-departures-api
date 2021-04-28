package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	v1 "github.com/Marchie/tf-experiment/lambda/internal/repository/redis/departures/metrolink/v1"
	mock_redis "github.com/Marchie/tf-experiment/lambda/pkg/mocks/redis"
	"github.com/aws/aws-sdk-go/aws"
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

func givenDeparturesForAtcoCode9400ZZMASTP1(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       0,
			Destination: "Rochdale via Oldham",
			Carriages:   "Single",
			Status:      "Arrived",
			Wait:        "0",
			Platform:    aws.String("D"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       1,
			Destination: "Victoria",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "5",
			Platform:    aws.String("D"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       2,
			Destination: "Rochdale via Oldham",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "12",
			Platform:    aws.String("D"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
	}
}

func givenDeparturesForAtcoCode9400ZZMASTP2(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       0,
			Destination: "Bury",
			Carriages:   "Single",
			Status:      "Departing",
			Wait:        "0",
			Platform:    aws.String("C"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       1,
			Destination: "Piccadilly",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "3",
			Platform:    aws.String("C"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       2,
			Destination: "Ashton-under-Lyne",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "7",
			Platform:    aws.String("C"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
	}
}

func givenDeparturesToStore(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       0,
			Destination: "Rochdale via Oldham",
			Carriages:   "Single",
			Status:      "Arrived",
			Wait:        "0",
			Platform:    aws.String("D"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       1,
			Destination: "Victoria",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "5",
			Platform:    aws.String("D"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       2,
			Destination: "Rochdale via Oldham",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "12",
			Platform:    aws.String("D"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       0,
			Destination: "Bury",
			Carriages:   "Single",
			Status:      "Departing",
			Wait:        "0",
			Platform:    aws.String("C"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       1,
			Destination: "Piccadilly",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "3",
			Platform:    aws.String("C"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       2,
			Destination: "Ashton-under-Lyne",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "7",
			Platform:    aws.String("C"),
			LastUpdated: time.Date(2021, time.March, 29, 11, 3, 40, 0, time.UTC),
		},
	}
}

func TestMetrolinkDeparturesRepository_Get(t *testing.T) {
	t.Run(`Given a populated Redis departures repository
When Get is called with an AtcoCode
Then departures for that AtcoCode are returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		atcoCode := "9400ZZMASTP1"

		expDepartures := givenDeparturesForAtcoCode9400ZZMASTP1(t)

		var departuresFromRedis bytes.Buffer
		if err := json.NewEncoder(&departuresFromRedis).Encode(expDepartures); err != nil {
			t.Fatal(err)
		}

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "departures", atcoCode)).Return(departuresFromRedis.Bytes(), nil),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		departures, err := metrolinkDeparturesRepository.Get(ctx, atcoCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, expDepartures, departures)
	})

	t.Run(`Given an error occurs getting a Redis connection from the connection pool
When Get is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresTimeToLive := time.Second * 15

		atcoCode := "9400ZZMASTP1"

		pool := mock_redis.NewMockPooler(ctrl)
		poolErr := errors.New("FUBAR")
		pool.EXPECT().GetContext(ctx).Return(nil, poolErr)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		departures, err := metrolinkDeparturesRepository.Get(ctx, atcoCode)

		// Then
		assert.Nil(t, departures)
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

		departuresKeyPrefix := "departures"

		atcoCode := "9400ZZMASTP1"

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)
		connErr := errors.New("FUBAR")

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "departures", atcoCode)).Return(nil, connErr),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		departures, err := metrolinkDeparturesRepository.Get(ctx, atcoCode)

		// Then
		assert.Nil(t, departures)
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

		departuresKeyPrefix := "departures"

		departuresTimeToLive := time.Second * 15

		atcoCode := "9400ZZMASTP1"

		conn := mock_redis.NewMockConn(ctrl)
		connCloseErr := errors.New("FUBAR")

		expDepartures := givenDeparturesForAtcoCode9400ZZMASTP1(t)

		var departuresFromRedis bytes.Buffer
		if err := json.NewEncoder(&departuresFromRedis).Encode(expDepartures); err != nil {
			t.Fatal(err)
		}

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "departures", atcoCode)).Return(departuresFromRedis.Bytes(), nil),
			conn.EXPECT().Close().Return(connCloseErr),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		departures, err := metrolinkDeparturesRepository.Get(ctx, atcoCode)

		// Then
		assert.NotNil(t, departures)
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

		departuresKeyPrefix := "departures"

		departuresTimeToLive := time.Second * 15

		atcoCode := "9400ZZMASTP1"

		conn := mock_redis.NewMockConn(ctrl)

		gomock.InOrder(
			conn.EXPECT().Do("GET", fmt.Sprintf("%s_%s", "departures", atcoCode)).Return([]byte("x"), nil),
			conn.EXPECT().Close().Return(nil),
		)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		departures, err := metrolinkDeparturesRepository.Get(ctx, atcoCode)

		// Then
		assert.Nil(t, departures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid character 'x' looking for beginning of value")
	})
}

func TestMetrolinkDeparturesRepository_Store(t *testing.T) {
	t.Run(`Given a slice of Metrolink departures
When Store is called
Then the departures are stored in Redis grouped by AtcoCode`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresToStore := givenDeparturesToStore(t)

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		var departuresFor9400ZZMASTP1 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP1).Encode(givenDeparturesForAtcoCode9400ZZMASTP1(t)); err != nil {
			t.Fatal(err)
		}

		var departuresFor9400ZZMASTP2 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP2).Encode(givenDeparturesForAtcoCode9400ZZMASTP2(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "departures_9400ZZMASTP1", departuresFor9400ZZMASTP1.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "departures_9400ZZMASTP2", departuresFor9400ZZMASTP2.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil).Times(2)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		err := metrolinkDeparturesRepository.Store(ctx, departuresToStore)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given an error occurs getting a Redis connection from the connection pool
When Store is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresToStore := givenDeparturesToStore(t)

		departuresTimeToLive := time.Second * 15

		pool := mock_redis.NewMockPooler(ctrl)
		poolErr := errors.New("FUBAR")
		pool.EXPECT().GetContext(ctx).Return(nil, poolErr)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		err := metrolinkDeparturesRepository.Store(ctx, departuresToStore)

		// Then
		assert.NotNil(t, err)
		assert.Equal(t, poolErr, err)
	})

	t.Run(`Given an error occurs sending data to Redis
When Store is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresToStore := givenDeparturesToStore(t)

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		var departuresFor9400ZZMASTP1 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP1).Encode(givenDeparturesForAtcoCode9400ZZMASTP1(t)); err != nil {
			t.Fatal(err)
		}

		var departuresFor9400ZZMASTP2 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP2).Encode(givenDeparturesForAtcoCode9400ZZMASTP2(t)); err != nil {
			t.Fatal(err)
		}

		connErr := errors.New("FUBAR")
		conn.EXPECT().Send("SET", "departures_9400ZZMASTP1", departuresFor9400ZZMASTP1.String(), "PX", int64(15000)).Return(connErr)
		conn.EXPECT().Send("SET", "departures_9400ZZMASTP2", departuresFor9400ZZMASTP2.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil).Times(1)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		err := metrolinkDeparturesRepository.Store(ctx, departuresToStore)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error sending/flushing Redis connection: 1 error occurred:\n\t* error sending Redis command for AtcoCode 9400ZZMASTP1: FUBAR\n\n")
	})

	t.Run(`Given an error occurs flushing the Redis connection
When Store is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresToStore := givenDeparturesToStore(t)

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		var departuresFor9400ZZMASTP1 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP1).Encode(givenDeparturesForAtcoCode9400ZZMASTP1(t)); err != nil {
			t.Fatal(err)
		}

		var departuresFor9400ZZMASTP2 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP2).Encode(givenDeparturesForAtcoCode9400ZZMASTP2(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "departures_9400ZZMASTP1", departuresFor9400ZZMASTP1.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "departures_9400ZZMASTP2", departuresFor9400ZZMASTP2.String(), "PX", int64(15000)).Return(nil)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Flush().Return(connErr)
		conn.EXPECT().Receive().Return("OK", nil).Times(2)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		err := metrolinkDeparturesRepository.Store(ctx, departuresToStore)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error sending/flushing Redis connection: 1 error occurred:\n\t* error flushing Redis connection: FUBAR\n\n")
	})

	t.Run(`Given an error occurs receiving the response from the Redis connection
When Store is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		departuresKeyPrefix := "departures"

		departuresToStore := givenDeparturesToStore(t)

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		var departuresFor9400ZZMASTP1 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP1).Encode(givenDeparturesForAtcoCode9400ZZMASTP1(t)); err != nil {
			t.Fatal(err)
		}

		var departuresFor9400ZZMASTP2 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP2).Encode(givenDeparturesForAtcoCode9400ZZMASTP2(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "departures_9400ZZMASTP1", departuresFor9400ZZMASTP1.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "departures_9400ZZMASTP2", departuresFor9400ZZMASTP2.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Receive().Return(nil, connErr)
		conn.EXPECT().Close().Return(nil)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		err := metrolinkDeparturesRepository.Store(ctx, departuresToStore)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error receiving on Redis connection: 1 error occurred:\n\t* FUBAR\n\n")
	})

	t.Run(`Given an error occurs returning the Redis connection to the pool
When Store is called
Then an error message is logged`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		departuresKeyPrefix := "departures"

		departuresToStore := givenDeparturesToStore(t)

		departuresTimeToLive := time.Second * 15

		conn := mock_redis.NewMockConn(ctrl)

		var departuresFor9400ZZMASTP1 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP1).Encode(givenDeparturesForAtcoCode9400ZZMASTP1(t)); err != nil {
			t.Fatal(err)
		}

		var departuresFor9400ZZMASTP2 bytes.Buffer
		if err := json.NewEncoder(&departuresFor9400ZZMASTP2).Encode(givenDeparturesForAtcoCode9400ZZMASTP2(t)); err != nil {
			t.Fatal(err)
		}

		conn.EXPECT().Send("SET", "departures_9400ZZMASTP1", departuresFor9400ZZMASTP1.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Send("SET", "departures_9400ZZMASTP2", departuresFor9400ZZMASTP2.String(), "PX", int64(15000)).Return(nil)
		conn.EXPECT().Flush().Return(nil)
		conn.EXPECT().Receive().Return("OK", nil).Times(2)
		connErr := errors.New("FUBAR")
		conn.EXPECT().Close().Return(connErr)

		pool := mock_redis.NewMockPooler(ctrl)
		pool.EXPECT().GetContext(ctx).Return(conn, nil)

		metrolinkDeparturesRepository := v1.NewMetrolinkDeparturesRepository(logger, pool, departuresKeyPrefix, departuresTimeToLive)

		// When
		err := metrolinkDeparturesRepository.Store(ctx, departuresToStore)

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
