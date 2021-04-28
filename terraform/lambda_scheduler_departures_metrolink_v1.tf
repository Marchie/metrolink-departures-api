resource "aws_lambda_function" "scheduler_departures_metrolink_v1_lambda_function" {
  depends_on = [
    aws_iam_role_policy_attachment.scheduler_departures_metrolink_v1_basic_execution_role_policy_attachment,
    aws_cloudwatch_log_group.scheduler_departures_metrolink_v1_cloudwatch_log_group,
    aws_cloudwatch_event_rule.every_one_minute,
    aws_sqs_queue.metrolink_departures_loader_v1_sqs_queue,
    aws_sqs_queue.metrolink_departures_loader_sqs_queue_deadletter,
  ]

  function_name = "scheduler_departures_metrolink_v1"

  s3_bucket = "marchie-build-artifacts"
  s3_key    = "v0.0.1/scheduler_departures_metrolink_v1.zip"

  handler = "scheduler_departures_metrolink_v1"
  runtime = "go1.x"

  timeout     = 3
  memory_size = 128

  reserved_concurrent_executions = 1

  role = aws_iam_role.scheduler_departures_metrolink_v1_iam_role.arn

  environment {
    variables = {
      FREQUENCY     = "3s"
      HORIZON       = "60s"
      LOG_LEVEL     = "0"
      SQS_QUEUE_URL = aws_sqs_queue.metrolink_departures_loader_v1_sqs_queue.id
    }
  }
}

resource "aws_iam_role" "scheduler_departures_metrolink_v1_iam_role" {
  name = "scheduler_departures_metrolink_v1"

  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role_iam_policy_document.json
}

resource "aws_iam_role_policy_attachment" "scheduler_departures_metrolink_v1_basic_execution_role_policy_attachment" {
  role       = aws_iam_role.scheduler_departures_metrolink_v1_iam_role.name
  policy_arn = data.aws_iam_policy.AWSLambdaBasicExecutionRole.arn
}

data "aws_iam_policy_document" "scheduler_departures_metrolink_v1_policy_document" {
  statement {
    actions = [
      "sqs:SendMessage"
    ]

    resources = [
      aws_sqs_queue.metrolink_departures_loader_v1_sqs_queue.arn
    ]

    effect = "Allow"
  }
}

resource "aws_iam_role_policy" "scheduler_departures_metrolink_v1_role_policy" {
  role   = aws_iam_role.scheduler_departures_metrolink_v1_iam_role.name
  policy = data.aws_iam_policy_document.scheduler_departures_metrolink_v1_policy_document.json
}

resource "aws_lambda_permission" "scheduler_departures_metrolink_v1_lambda_permission" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.scheduler_departures_metrolink_v1_lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_one_minute.arn
}

resource "aws_cloudwatch_log_group" "scheduler_departures_metrolink_v1_cloudwatch_log_group" {
  name              = "/aws/lambda/scheduler_departures_metrolink_v1"
  retention_in_days = var.cloudwatch_logs_retention_in_days
}

resource "aws_cloudwatch_event_target" "scheduler_invoke_every_one_minute" {
  depends_on = [
    aws_cloudwatch_event_rule.every_one_minute
  ]

  rule      = aws_cloudwatch_event_rule.every_one_minute.name
  target_id = "scheduler_departures_metrolink_v1_lambda"
  arn       = aws_lambda_function.scheduler_departures_metrolink_v1_lambda_function.arn
}

resource "aws_cloudwatch_metric_alarm" "scheduler_departures_metrolink_v1_alarm_on_error_and_insufficient_data" {
  alarm_name          = "scheduler-departures-metrolink-v1-alarm-if-error"
  alarm_description   = "Error occurred scheduling Metrolink departures data loader"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  period              = 120
  namespace           = "AWS/Lambda"
  metric_name         = "Errors"
  dimensions = {
    FunctionName = aws_lambda_function.scheduler_departures_metrolink_v1_lambda_function.function_name
  }
  statistic       = "Sum"
  threshold       = 1
  actions_enabled = true
  alarm_actions = [
    aws_sns_topic.production_notifications.arn
  ]
  ok_actions = [
    aws_sns_topic.production_notifications.arn
  ]
  insufficient_data_actions = [
    aws_sns_topic.production_notifications.arn
  ]
}
