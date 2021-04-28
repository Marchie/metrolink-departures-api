package main

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core/departures/metrolink/api"
	"github.com/Marchie/tf-experiment/lambda/internal/logger"
	v1 "github.com/Marchie/tf-experiment/lambda/internal/repository/redis/departures/metrolink/v1"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/redis/naptan"
	v12 "github.com/Marchie/tf-experiment/lambda/internal/repository/redis/system/status/v1"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/apigw"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/plaid/go-envvar/envvar"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	LogLevel                                           int8          `envvar:"LOG_LEVEL" default:"0"`
	MetrolinkDeparturesStaleDataThreshold              time.Duration `envvar:"METROLINK_DEPARTURES_STALE_DATA_THRESHOLD" default:"30s"`
	RedisMetrolinkDeparturesServerAddress              string        `envvar:"REDIS_METROLINK_DEPARTURES_SERVER_ADDRESS"`
	RedisMetrolinkDeparturesKeyPrefix                  string        `envvar:"REDIS_METROLINK_DEPARTURES_KEY_PREFIX" default:"metrolink_departures"`
	RedisMetrolinkDeparturesServiceStatusServerAddress string        `envvar:"REDIS_METROLINK_DEPARTURES_SERVICE_STATUS_SERVER_ADDRESS"`
	RedisMetrolinkDeparturesServiceStatusKey           string        `envvar:"REDIS_METROLINK_DEPARTURES_SERVICE_STATUS_KEY" default:"metrolink_departures_service_status"`
	RedisStopsInAreaServerAddress                      string        `envvar:"REDIS_STOPS_IN_AREA_SERVER_ADDRESS"`
	RedisStopsInAreaKeyPrefix                          string        `envvar:"REDIS_STOPS_IN_AREA_KEY_PREFIX" default:"stops_in_area"`
	StopAreaCodeOrAtcoCodeApiGatewayPathParameter      string        `envvar:"STOP_AREA_CODE_OR_ATCO_CODE_API_GATEWAY_PATH_PARAMETER"`
	TimeLocation                                       string        `envvar:"TIME_LOCATION" default:"Europe/London"`
}

func main() {
	var cfg Config
	if err := envvar.Parse(&cfg); err != nil {
		panic(errors.Wrap(err, "error parsing config"))
	}

	baseLogger, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(errors.Wrap(err, "error creating new Logger"))
	}

	metrolinkDeparturesPool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", cfg.RedisMetrolinkDeparturesServerAddress)
		},
	}

	metrolinkDeparturesSystemStatusPool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", cfg.RedisMetrolinkDeparturesServiceStatusServerAddress)
		},
	}

	stopsInAreaPool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", cfg.RedisStopsInAreaServerAddress)
		},
	}

	timeLocation, err := time.LoadLocation(cfg.TimeLocation)
	if err != nil {
		panic(errors.Wrap(err, "error loading time location"))
	}

	lambda.Start(func(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		lc, _ := lambdacontext.FromContext(ctx)

		childLogger := baseLogger.With(zap.String("awsRequestId", lc.AwsRequestID))

		stopsInAreaGetter := naptan.NewNaptanRedis(childLogger, stopsInAreaPool, cfg.RedisStopsInAreaKeyPrefix, 0)

		metrolinkDeparturesGetter := v1.NewMetrolinkDeparturesRepository(childLogger, metrolinkDeparturesPool, cfg.RedisMetrolinkDeparturesKeyPrefix, 0)

		systemStatusGetter := v12.NewMetrolinkDeparturesSystemStatusRepository(childLogger, metrolinkDeparturesSystemStatusPool, cfg.RedisMetrolinkDeparturesServiceStatusKey)

		metrolinkDeparturesApi := api.NewApi(childLogger, stopsInAreaGetter, metrolinkDeparturesGetter, systemStatusGetter, time.Now, cfg.MetrolinkDeparturesStaleDataThreshold, timeLocation)

		return apigw.NewMetrolinkDeparturesAwsApiGateway(childLogger, metrolinkDeparturesApi, cfg.StopAreaCodeOrAtcoCodeApiGatewayPathParameter).Handler(ctx, event)
	})
}
