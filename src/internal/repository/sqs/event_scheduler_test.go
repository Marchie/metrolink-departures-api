package sqs

import (
	"context"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	mock_sqsiface "github.com/Marchie/tf-experiment/lambda/pkg/mocks/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func given11Events(t *testing.T, now time.Time, payload string) []*domain.Event {
	t.Helper()

	return []*domain.Event{
		{
			StartTime: now,
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 10),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 20),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 30),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 40),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 50),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 60),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 70),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 80),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 90),
			Payload:   payload,
		},
		{
			StartTime: now.Add(time.Second * 100),
			Payload:   payload,
		},
	}
}

func given11EventsMessageBatch1Expectation(t *testing.T, messageIdPrefix string, now time.Time, payload string, sqsQueueUrl string) *sqs.SendMessageBatchInput {
	t.Helper()

	return &sqs.SendMessageBatchInput{
		Entries: []*sqs.SendMessageBatchRequestEntry{
			{
				DelaySeconds: aws.Int64(0),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(10),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(10*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(20),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(20*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(30),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(30*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(40),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(40*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(50),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(50*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(60),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(60*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(70),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(70*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(80),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(80*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
			{
				DelaySeconds: aws.Int64(90),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(90*time.Second).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
		},
		QueueUrl: &sqsQueueUrl,
	}
}

func given11EventsMessageBatch2Expectation(t *testing.T, messageIdPrefix string, now time.Time, payload string, sqsQueueUrl string) *sqs.SendMessageBatchInput {
	t.Helper()

	return &sqs.SendMessageBatchInput{
		Entries: []*sqs.SendMessageBatchRequestEntry{
			{
				DelaySeconds: aws.Int64(100),
				Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(time.Second*100).Format("2006-01-02T15_04_05"))),
				MessageBody:  &payload,
			},
		},
		QueueUrl: &sqsQueueUrl,
	}
}

func TestEventScheduler_Schedule(t *testing.T) {
	t.Run(`Given a batch of 11 events
When Schedule is called
Then the events are sent as messages to SQS in batches with a maximum size of 10
And the messages have an appropriate delay`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		sqsClient := mock_sqsiface.NewMockSQSAPI(ctrl)

		sqsQueueUrl := "http://sqs.queue"

		messageIdPrefix := "foo"

		now := time.Date(2021, time.March, 28, 16, 54, 15, 0, time.UTC)

		payload := "bar"

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := NewEventScheduler(logger, sqsClient, sqsQueueUrl, messageIdPrefix, currentTimeFunc)

		events := given11Events(t, now, payload)

		expMessageBatch1 := given11EventsMessageBatch1Expectation(t, messageIdPrefix, now, payload, sqsQueueUrl)

		messageBatch1Output := &sqs.SendMessageBatchOutput{
			Failed:     nil,
			Successful: nil,
		}

		expMessageBatch2 := given11EventsMessageBatch2Expectation(t, messageIdPrefix, now, payload, sqsQueueUrl)

		messageBatch2Output := &sqs.SendMessageBatchOutput{
			Failed:     nil,
			Successful: nil,
		}

		gomock.InOrder(
			sqsClient.EXPECT().SendMessageBatchWithContext(ctx, expMessageBatch1).Return(messageBatch1Output, nil),
			sqsClient.EXPECT().SendMessageBatchWithContext(ctx, expMessageBatch2).Return(messageBatch2Output, nil),
		)

		// When
		err := eventScheduler.Schedule(ctx, events)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given a batch of events
When Schedule is called
And messages are queued successfully
Then information is logged at the debug level`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		sqsClient := mock_sqsiface.NewMockSQSAPI(ctrl)

		sqsQueueUrl := "http://sqs.queue"

		messageIdPrefix := "foo"

		now := time.Date(2021, time.March, 28, 16, 54, 15, 0, time.UTC)

		payload := "bar"

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := NewEventScheduler(logger, sqsClient, sqsQueueUrl, messageIdPrefix, currentTimeFunc)

		events := []*domain.Event{
			{
				StartTime: now,
				Payload:   payload,
			},
		}

		id := fmt.Sprintf("%s_%s", messageIdPrefix, now.Format("2006-01-02T15_04_05"))
		msgId := "abc123"

		messageBatchOutput := &sqs.SendMessageBatchOutput{
			Successful: []*sqs.SendMessageBatchResultEntry{
				{
					Id:        &id,
					MessageId: &msgId,
				},
			},
		}

		sqsClient.EXPECT().SendMessageBatchWithContext(ctx, gomock.Any()).Return(messageBatchOutput, nil)

		// When
		err := eventScheduler.Schedule(ctx, events)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, 1, observedLogs.Len())

		loggedRecords := observedLogs.TakeAll()
		assert.Equal(t, zapcore.DebugLevel, loggedRecords[0].Level)
		assert.Equal(t, "successfully enqueued SQS message", loggedRecords[0].Message)
		assert.Equal(t, "id", loggedRecords[0].Context[0].Key)
		assert.Equal(t, id, loggedRecords[0].Context[0].String)
		assert.Equal(t, "messageId", loggedRecords[0].Context[1].Key)
		assert.Equal(t, msgId, loggedRecords[0].Context[1].String)
	})

	t.Run(`Given a batch of events
When Schedule is called
And messages fail to be queued
Then information is logged at the error level
And an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zapcore.DebugLevel)
		logger := zap.New(zapCore)

		sqsClient := mock_sqsiface.NewMockSQSAPI(ctrl)

		sqsQueueUrl := "http://sqs.queue"

		messageIdPrefix := "foo"

		now := time.Date(2021, time.March, 28, 16, 54, 15, 0, time.UTC)

		payload := "bar"

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := NewEventScheduler(logger, sqsClient, sqsQueueUrl, messageIdPrefix, currentTimeFunc)

		events := []*domain.Event{
			{
				StartTime: now,
				Payload:   payload,
			},
		}

		id := fmt.Sprintf("%s_%s", messageIdPrefix, now.Format("2006-01-02T15_04_05"))
		code := "abc123"
		msg := "FUBAR"
		senderFault := true

		messageBatchOutput := &sqs.SendMessageBatchOutput{
			Failed: []*sqs.BatchResultErrorEntry{
				{
					Code:        &code,
					Id:          &id,
					Message:     &msg,
					SenderFault: &senderFault,
				},
			},
		}

		sqsClient.EXPECT().SendMessageBatchWithContext(ctx, gomock.Any()).Return(messageBatchOutput, nil)

		// When
		err := eventScheduler.Schedule(ctx, events)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "1 error occurred:\n\t* error enqueuing 1 message(s)\n\n")
		assert.Equal(t, 1, observedLogs.Len())

		loggedRecords := observedLogs.TakeAll()
		assert.Equal(t, zapcore.ErrorLevel, loggedRecords[0].Level)
		assert.Equal(t, "failed to enqueue SQS message", loggedRecords[0].Message)
		assert.Equal(t, "id", loggedRecords[0].Context[0].Key)
		assert.Equal(t, id, loggedRecords[0].Context[0].String)
		assert.Equal(t, "message", loggedRecords[0].Context[1].Key)
		assert.Equal(t, msg, loggedRecords[0].Context[1].String)
		assert.Equal(t, "code", loggedRecords[0].Context[2].Key)
		assert.Equal(t, code, loggedRecords[0].Context[2].String)
		assert.Equal(t, "senderFault", loggedRecords[0].Context[3].Key)
		assert.Equal(t, int64(1), loggedRecords[0].Context[3].Integer)
	})

	t.Run(`Given an event in a batch has start time that is more than 900 seconds into the future
When Schedule is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		sqsClient := mock_sqsiface.NewMockSQSAPI(ctrl)

		sqsQueueUrl := "http://sqs.queue"

		messageIdPrefix := "foo"

		now := time.Date(2021, time.March, 28, 16, 54, 15, 0, time.UTC)

		payload := "bar"

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := NewEventScheduler(logger, sqsClient, sqsQueueUrl, messageIdPrefix, currentTimeFunc)

		events := []*domain.Event{
			{
				StartTime: now.Add(time.Second * 901),
				Payload:   payload,
			},
		}

		// When
		err := eventScheduler.Schedule(ctx, events)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "1 error occurred:\n\t* event start time 2021-03-28 17:09:16 +0000 UTC is too far in the future (901 seconds); SQS items can only be delayed by up to 900 seconds\n\n")
	})

	t.Run(`Given an event in a batch has start time that is in the past
When Schedule is called
Then the event start time is normalised to have a delay time of 0`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		sqsClient := mock_sqsiface.NewMockSQSAPI(ctrl)

		sqsQueueUrl := "http://sqs.queue"

		messageIdPrefix := "foo"

		now := time.Date(2021, time.March, 28, 16, 54, 15, 0, time.UTC)

		payload := "bar"

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := NewEventScheduler(logger, sqsClient, sqsQueueUrl, messageIdPrefix, currentTimeFunc)

		events := []*domain.Event{
			{
				StartTime: now.Add(-time.Second),
				Payload:   payload,
			},
		}

		expMessageBatch := &sqs.SendMessageBatchInput{
			Entries: []*sqs.SendMessageBatchRequestEntry{
				{
					DelaySeconds: aws.Int64(0),
					Id:           aws.String(fmt.Sprintf("%s_%s", messageIdPrefix, now.Add(-time.Second).Format("2006-01-02T15_04_05"))),
					MessageBody:  &payload,
				},
			},
			QueueUrl: &sqsQueueUrl,
		}

		messageBatchOutput := &sqs.SendMessageBatchOutput{
			Failed:     nil,
			Successful: nil,
		}

		sqsClient.EXPECT().SendMessageBatchWithContext(ctx, expMessageBatch).Return(messageBatchOutput, nil)

		// When
		err := eventScheduler.Schedule(ctx, events)

		// Then
		assert.Nil(t, err)
	})

	t.Run(`Given the SQS client returns an error
When Schedule is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		sqsClient := mock_sqsiface.NewMockSQSAPI(ctrl)

		sqsQueueUrl := "http://sqs.queue"

		messageIdPrefix := "foo"

		now := time.Date(2021, time.March, 28, 16, 54, 15, 0, time.UTC)

		payload := "bar"

		currentTimeFunc := func() time.Time {
			return now
		}

		eventScheduler := NewEventScheduler(logger, sqsClient, sqsQueueUrl, messageIdPrefix, currentTimeFunc)

		events := []*domain.Event{
			{
				StartTime: now,
				Payload:   payload,
			},
		}

		sqsClientErr := errors.New("FUBAR")
		sqsClient.EXPECT().SendMessageBatchWithContext(ctx, gomock.Any()).Return(nil, sqsClientErr)

		// When
		err := eventScheduler.Schedule(ctx, events)

		// Then
		assert.NotNil(t, err)
		assert.EqualError(t, err, "1 error occurred:\n\t* FUBAR\n\n")
	})
}
