package main

import (
	"context"
	v1 "github.com/Marchie/tf-experiment/lambda/internal/core/schedule/departures/metrolink/v1"
	"github.com/Marchie/tf-experiment/lambda/internal/logger"
	sqs2 "github.com/Marchie/tf-experiment/lambda/internal/repository/sqs"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/cloudwatch"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/plaid/go-envvar/envvar"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Frequency          time.Duration `envvar:"FREQUENCY"`
	Horizon            time.Duration `envvar:"HORIZON"`
	LogLevel           int8          `envvar:"LOG_LEVEL"`
	SQSMessageIdPrefix string        `envvar:"SQS_MESSAGE_ID_PREFIX" default:"msg"`
	SQSQueueUrl        string        `envvar:"SQS_QUEUE_URL"`
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

	sess, err := session.NewSession()
	if err != nil {
		panic(errors.Wrap(err, "error creating AWS Session"))
	}

	sqsClient := sqs.New(sess)

	lambda.Start(func(ctx context.Context) error {
		lc, _ := lambdacontext.FromContext(ctx)

		childLogger := baseLogger.With(zap.String("awsRequestId", lc.AwsRequestID))

		sqsEventScheduler := sqs2.NewEventScheduler(childLogger, sqsClient, cfg.SQSQueueUrl, cfg.SQSMessageIdPrefix, time.Now)

		metrolinkDataLoadScheduler, err := v1.NewMetrolinkDeparturesDataLoadScheduler(childLogger, sqsEventScheduler, cfg.Horizon, cfg.Frequency, "x", time.Now)
		if err != nil {
			return err
		}

		return cloudwatch.NewMetrolinkDeparturesDataLoadScheduler(childLogger, metrolinkDataLoadScheduler).Handler(ctx)
	})
}
