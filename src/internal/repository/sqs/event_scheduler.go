package sqs

import (
	"context"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	"time"
)

type EventScheduler struct {
	logger           *zap.Logger
	sqsClient        sqsiface.SQSAPI
	sqsQueueUrl      string
	messageIdPrefix  string
	currentTimeFunc  func() time.Time
	maximumBatchSize int
}

func NewEventScheduler(logger *zap.Logger, sqsClient sqsiface.SQSAPI, sqsQueueUrl string, messageIdPrefix string, currentTimeFunc func() time.Time) *EventScheduler {
	return &EventScheduler{
		logger:           logger,
		sqsClient:        sqsClient,
		sqsQueueUrl:      sqsQueueUrl,
		messageIdPrefix:  messageIdPrefix,
		currentTimeFunc:  currentTimeFunc,
		maximumBatchSize: 10,
	}
}

func (e *EventScheduler) Schedule(ctx context.Context, events []*domain.Event) error {
	batchedEvents := e.splitEventsIntoBatches(events)

	var errs error

	for _, batch := range batchedEvents {
		sendMessageBatchRequestEntries, err := e.convertBatchToSendMessageBatchRequestEntry(batch)
		if err != nil {
			errs = multierror.Append(errs, err)
		}

		if len(sendMessageBatchRequestEntries) == 0 {
			continue
		}

		sqsSendMessageBatchOutput, err := e.sqsClient.SendMessageBatchWithContext(ctx, &sqs.SendMessageBatchInput{
			Entries:  sendMessageBatchRequestEntries,
			QueueUrl: &e.sqsQueueUrl,
		})
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if err := e.handleSQSResponse(sqsSendMessageBatchOutput); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (e *EventScheduler) convertBatchToSendMessageBatchRequestEntry(events []*domain.Event) ([]*sqs.SendMessageBatchRequestEntry, error) {
	var sendMessageBatchRequestEntries []*sqs.SendMessageBatchRequestEntry

	var errs error

	for _, event := range events {
		sendMessageBatchRequestEntry, err := e.convertToSendMessageBatchRequestEntry(event)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		sendMessageBatchRequestEntries = append(sendMessageBatchRequestEntries, sendMessageBatchRequestEntry)
	}

	return sendMessageBatchRequestEntries, errs
}

func (e *EventScheduler) convertToSendMessageBatchRequestEntry(event *domain.Event) (*sqs.SendMessageBatchRequestEntry, error) {
	delaySeconds := int64(event.StartTime.Sub(e.currentTimeFunc()).Seconds())

	if delaySeconds < 0 {
		delaySeconds = 0
	}

	if delaySeconds > 900 {
		return nil, fmt.Errorf("event start time %s is too far in the future (%d seconds); SQS items can only be delayed by up to 900 seconds", event.StartTime, delaySeconds)
	}

	id := fmt.Sprintf("%s_%s", e.messageIdPrefix, event.StartTime.Format("2006-01-02T15_04_05"))

	return &sqs.SendMessageBatchRequestEntry{
		DelaySeconds: &delaySeconds,
		Id:           &id,
		MessageBody:  &event.Payload,
	}, nil
}

func (e *EventScheduler) handleSQSResponse(sqsSendMessageBatchOutput *sqs.SendMessageBatchOutput) error {
	failedCount := 0

	for _, failed := range sqsSendMessageBatchOutput.Failed {
		e.logger.Error("failed to enqueue SQS message", zap.Stringp("id", failed.Id), zap.Stringp("message", failed.Message), zap.Stringp("code", failed.Code), zap.Boolp("senderFault", failed.SenderFault))

		failedCount++
	}

	for _, successful := range sqsSendMessageBatchOutput.Successful {
		e.logger.Debug("successfully enqueued SQS message", zap.Stringp("id", successful.Id), zap.Stringp("messageId", successful.MessageId))
	}

	if failedCount > 0 {
		return fmt.Errorf("error enqueuing %d message(s)", failedCount)
	}

	return nil
}

func (e *EventScheduler) splitEventsIntoBatches(events []*domain.Event) [][]*domain.Event {
	var batches [][]*domain.Event

	for i := 0; i < len(events); i += e.maximumBatchSize {
		end := i + e.maximumBatchSize
		if end > len(events) {
			end = len(events)
		}

		batches = append(batches, events[i:end])
	}

	return batches
}
