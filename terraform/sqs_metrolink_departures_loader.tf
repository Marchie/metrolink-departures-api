resource "aws_sqs_queue" "metrolink_departures_loader_v1_sqs_queue" {
  name                       = "metrolink-departures-loader-v1-sqs-queue"
  visibility_timeout_seconds = 15
  message_retention_seconds  = 60
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.metrolink_departures_loader_sqs_queue_deadletter.arn
    maxReceiveCount     = 1
  })
}

resource "aws_sqs_queue" "metrolink_departures_loader_sqs_queue_deadletter" {
  name = "metrolink-departures-loader-sqs-queue-deadletter"
}
