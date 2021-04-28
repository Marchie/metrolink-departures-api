package main

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core/naptan/loader"
	"github.com/Marchie/tf-experiment/lambda/internal/logger"
	naptan2 "github.com/Marchie/tf-experiment/lambda/internal/repository/api/http/dft/naptan"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/compression"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/redis/naptan"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/cloudwatch"
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
	HttpClientTimeout                        time.Duration `envvar:"HTTP_CLIENT_TIMEOUT" default:"15s"`
	LogLevel                                 int8          `envvar:"LOG_LEVEL" default:"0"`
	NaptanCsvUrl                             string        `envvar:"NAPTAN_CSV_URL"`
	NaptanStopsInAreaFilename                string        `envvar:"NAPTAN_STOPS_IN_AREA_FILENAME" default:"StopsInArea.csv"`
	NaptanStopsInAreaStopAreaCodeColumnIndex int           `envvar:"NAPTAN_STOPS_IN_AREA_STOP_AREA_CODE_COLUMN_INDEX" default:"0"`
	NaptanStopsInAreaAtcoCodeColumnIndex     int           `envvar:"NAPTAN_STOPS_IN_AREA_ATCO_CODE_COLUMN_INDEX" default:"1"`
	RedisServerAddress                       string        `envvar:"REDIS_SERVER_ADDRESS"`
	RedisStopsInAreaKeyPrefix                string        `envvar:"REDIS_STOPS_IN_AREA_KEY_PREFIX" default:"stops_in_area"`
	RedisStopsInAreaTimeToLive               time.Duration `envvar:"REDIS_STOPS_IN_AREA_TIME_TO_LIVE" default:"25h"`
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

	pool := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", cfg.RedisServerAddress)
		},
	}

	lambda.Start(func(ctx context.Context) error {
		lc, _ := lambdacontext.FromContext(ctx)

		childLogger := baseLogger.With(zap.String("awsRequestId", lc.AwsRequestID))

		httpZipFileFetcher := naptan2.NewRepository(childLogger, httpClient, cfg.NaptanCsvUrl)

		zipFileExtractor := compression.NewZipFileExtractor(childLogger)

		httpStopsInAreaFetcher := naptan2.NewCSV(childLogger, httpZipFileFetcher, zipFileExtractor, cfg.NaptanStopsInAreaFilename, cfg.NaptanStopsInAreaStopAreaCodeColumnIndex, cfg.NaptanStopsInAreaAtcoCodeColumnIndex)

		redisStopsInAreaStorer := naptan.NewNaptanRedis(childLogger, pool, cfg.RedisStopsInAreaKeyPrefix, cfg.RedisStopsInAreaTimeToLive)

		stopsInAreaLoader := loader.NewStopsInAreaLoader(childLogger, httpStopsInAreaFetcher, redisStopsInAreaStorer)

		naptanDataLoader := cloudwatch.NewNaptanDataLoader(childLogger, stopsInAreaLoader)

		return naptanDataLoader.Handler(ctx)
	})
}
