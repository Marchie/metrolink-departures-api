resource "aws_lambda_function" "dataloader_departures_metrolink_v1_lambda_function" {
  depends_on = [
    aws_iam_role_policy_attachment.dataloader_departures_metrolink_v1_iam_role_policy_attachment,
    aws_cloudwatch_log_group.dataloader_departures_metrolink_v1_cloudwatch_log_group,
    aws_elasticache_cluster.tfgm_com,
    aws_cloudwatch_event_rule.every_one_minute,
    aws_sqs_queue.metrolink_departures_loader_v1_sqs_queue,
    aws_sqs_queue.metrolink_departures_loader_sqs_queue_deadletter,
  ]

  function_name = "dataloader_departures_metrolink_v1"

  s3_bucket = "marchie-build-artifacts"
  s3_key    = "v0.0.1/dataloader_departures_metrolink_v1.zip"

  handler = "dataloader_departures_metrolink_v1"
  runtime = "go1.x"

  timeout     = 3
  memory_size = 128

  reserved_concurrent_executions = 1

  role = aws_iam_role.dataloader_departures_metrolink_v1_iam_role.arn

  vpc_config {
    security_group_ids = [
      aws_security_group.production_lambda_sg.id
    ]
    subnet_ids = [
      aws_subnet.production_private.id
    ]
  }

  environment {
    variables = {
      LOG_LEVEL                                                = "-1"
      METROLINK_DEPARTURES_STALE_DATA_THRESHOLD                = "30s"
      REDIS_METROLINK_DEPARTURES_SERVER_ADDRESS                = "${aws_elasticache_cluster.tfgm_com.cache_nodes.0.address}:${aws_elasticache_cluster.tfgm_com.cache_nodes.0.port}"
      REDIS_METROLINK_DEPARTURES_SERVICE_STATUS_SERVER_ADDRESS = "${aws_elasticache_cluster.tfgm_com.cache_nodes.0.address}:${aws_elasticache_cluster.tfgm_com.cache_nodes.0.port}"
      REDIS_METROLINK_DEPARTURES_TIME_TO_LIVE                  = "15s"
      TFGM_METROLINKS_API_KEY                                  = var.departures_metrolink_v1_tfgm_developer_api_key
      TFGM_METROLINKS_API_URL                                  = var.departures_metrolink_v1_tfgm_developer_metrolinks_url
    }
  }
}

resource "aws_iam_role" "dataloader_departures_metrolink_v1_iam_role" {
  name = "dataloader_departures_metrolink_v1"

  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role_iam_policy_document.json
}

resource "aws_iam_role_policy_attachment" "dataloader_departures_metrolink_v1_iam_role_policy_attachment" {
  role       = aws_iam_role.dataloader_departures_metrolink_v1_iam_role.name
  policy_arn = data.aws_iam_policy.AWSLambdaVPCAccessExecutionRole.arn
}

data "aws_iam_policy_document" "dataloader_departures_metrolink_v1_iam_policy_document" {
  statement {
    effect = "Allow"

    resources = [
      aws_sqs_queue.metrolink_departures_loader_v1_sqs_queue.arn
    ]

    actions = [
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
      "sqs:ReceiveMessage"
    ]
  }
}

resource "aws_iam_role_policy" "dataloader_departures_metrolink_v1_iam_role_policy" {
  policy = data.aws_iam_policy_document.dataloader_departures_metrolink_v1_iam_policy_document.json
  role   = aws_iam_role.dataloader_departures_metrolink_v1_iam_role.name
}

resource "aws_lambda_event_source_mapping" "dataloader_departures_metrolink_v1_sqs_event_source_mapping" {
  batch_size       = 1
  enabled          = true
  event_source_arn = aws_sqs_queue.metrolink_departures_loader_v1_sqs_queue.arn
  function_name    = aws_lambda_function.dataloader_departures_metrolink_v1_lambda_function.arn
}

resource "aws_cloudwatch_log_group" "dataloader_departures_metrolink_v1_cloudwatch_log_group" {
  name              = "/aws/lambda/dataloader_departures_metrolink_v1"
  retention_in_days = var.cloudwatch_logs_retention_in_days
}

resource "aws_cloudwatch_metric_alarm" "dataloader_departures_metrolink_v1_alarm_on_error_and_insufficient_data" {
  alarm_name          = "dataloader-departures-metrolink-v1-alarm-if-error"
  alarm_description   = "Error occurred loading Metrolink departures data"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  period              = 60
  namespace           = "AWS/Lambda"
  metric_name         = "Errors"
  dimensions = {
    FunctionName = aws_lambda_function.dataloader_departures_metrolink_v1_lambda_function.function_name
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

resource "aws_cloudwatch_metric_alarm" "dataloader_departures_metrolink_v1_alarm_on_low_invocations_count" {
  alarm_name          = "dataloader-departures-metrolink-v1-alarm-on-low-invocations"
  alarm_description   = "Metrolink departures data not retrieved with sufficient frequency"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 1
  period              = 300
  namespace           = "AWS/Lambda"
  metric_name         = "Invocations"
  dimensions = {
    FunctionName = aws_lambda_function.dataloader_departures_metrolink_v1_lambda_function.function_name
  }
  statistic       = "Sum"
  threshold       = 90
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
