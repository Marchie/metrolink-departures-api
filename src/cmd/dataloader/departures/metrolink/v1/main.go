package main

import (
	"context"
	loader2 "github.com/Marchie/tf-experiment/lambda/internal/core/departures/metrolink/loader"
	"github.com/Marchie/tf-experiment/lambda/internal/logger"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/api/http/tfgm/developer"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/filesystem"
	v1 "github.com/Marchie/tf-experiment/lambda/internal/repository/redis/departures/metrolink/v1"
	v12 "github.com/Marchie/tf-experiment/lambda/internal/repository/redis/system/status/v1"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/sqs"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/plaid/go-envvar/envvar"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Config struct {
	HttpClientTimeout                                  time.Duration `envvar:"HTTP_CLIENT_TIMEOUT" default:"1500ms"`
	LogLevel                                           int8          `envvar:"LOG_LEVEL" default:"0"`
	MetrolinkDeparturesStaleDataThreshold              time.Duration `envvar:"METROLINK_DEPARTURES_STALE_DATA_THRESHOLD" default:"30s"`
	RedisMetrolinkDeparturesServerAddress              string        `envvar:"REDIS_METROLINK_DEPARTURES_SERVER_ADDRESS"`
	RedisMetrolinkDeparturesKeyPrefix                  string        `envvar:"REDIS_METROLINK_DEPARTURES_KEY_PREFIX" default:"metrolink_departures"`
	RedisMetrolinkDeparturesTimeToLive                 time.Duration `envvar:"REDIS_METROLINK_DEPARTURES_TIME_TO_LIVE" default:"15s"`
	RedisMetrolinkDeparturesServiceStatusServerAddress string        `envvar:"REDIS_METROLINK_DEPARTURES_SERVICE_STATUS_SERVER_ADDRESS"`
	RedisMetrolinkDeparturesServiceStatusKey           string        `envvar:"REDIS_METROLINK_DEPARTURES_SERVICE_STATUS_KEY" default:"metrolink_departures_service_status"`
	TfgmMetrolinksApiKey                               string        `envvar:"TFGM_METROLINKS_API_KEY"`
	TfgmMetrolinksApiUrl                               string        `envvar:"TFGM_METROLINKS_API_URL"`
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

	httpClient := &http.Client{
		Timeout: cfg.HttpClientTimeout,
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

	lambda.Start(func(ctx context.Context) error {
		lc, _ := lambdacontext.FromContext(ctx)

		childLogger := baseLogger.With(zap.String("awsRequestId", lc.AwsRequestID))

		metrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(childLogger, httpClient, cfg.TfgmMetrolinksApiUrl, cfg.TfgmMetrolinksApiKey)

		platformNamer := filesystem.NewPlatformNamer(childLogger)

		redisMetrolinkDeparturesStorer := v1.NewMetrolinkDeparturesRepository(childLogger, metrolinkDeparturesPool, cfg.RedisMetrolinkDeparturesKeyPrefix, cfg.RedisMetrolinkDeparturesTimeToLive)

		redisMetrolinkDeparturesSystemStatusStorer := v12.NewMetrolinkDeparturesSystemStatusRepository(childLogger, metrolinkDeparturesSystemStatusPool, cfg.RedisMetrolinkDeparturesServiceStatusKey)

		metrolinkDeparturesLoader := loader2.NewMetrolinkDeparturesLoader(childLogger, metrolinkDataSource, platformNamer, redisMetrolinkDeparturesStorer, redisMetrolinkDeparturesSystemStatusStorer, time.Now, cfg.MetrolinkDeparturesStaleDataThreshold)

		metrolinkDeparturesDataLoader := sqs.NewMetrolinkDeparturesDataLoader(childLogger, metrolinkDeparturesLoader)

		return metrolinkDeparturesDataLoader.Handler(ctx)
	})
}
