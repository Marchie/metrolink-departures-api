# scheduler-departures-metrolink-v1

A Lambda function which creates messages on an SQS queue, with incrementing delays. The messages on the queue are used
to invoke the [dataloader-departures-metrolink-v1 Lambda function](../../../../dataloader/departures/metrolink/v1/README.md)
with a sub-minute frequency.
