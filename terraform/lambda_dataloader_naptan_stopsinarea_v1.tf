resource "aws_lambda_function" "dataloader_naptan_stopsinarea_v1_lambda_function" {
  depends_on = [
    aws_iam_role_policy_attachment.dataloader_naptan_stopsinarea_v1_iam_role_policy_attachment,
    aws_cloudwatch_log_group.dataloader_naptan_stopsinarea_v1_cloudwatch_log_group,
    aws_elasticache_cluster.tfgm_com,
    aws_cloudwatch_event_rule.every_day_at_0330,
  ]

  function_name = "dataloader_naptan_stopsinarea_v1"

  s3_bucket = "marchie-build-artifacts"
  s3_key    = "v0.0.1/dataloader_naptan_stopsinarea_v1.zip"

  handler = "dataloader_naptan_stopsinarea_v1"
  runtime = "go1.x"

  timeout     = 30
  memory_size = 1024

  reserved_concurrent_executions = 1

  role = aws_iam_role.dataloader_naptan_stopsinarea_v1_iam_role.arn

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
      LOG_LEVEL            = "-1"
      NAPTAN_CSV_URL       = var.naptan_csv_zip_url
      REDIS_SERVER_ADDRESS = "${aws_elasticache_cluster.tfgm_com.cache_nodes.0.address}:${aws_elasticache_cluster.tfgm_com.cache_nodes.0.port}"
    }
  }
}

resource "aws_iam_role" "dataloader_naptan_stopsinarea_v1_iam_role" {
  name = "dataloader_naptan_stopsinarea_v1"

  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role_iam_policy_document.json
}

resource "aws_iam_role_policy_attachment" "dataloader_naptan_stopsinarea_v1_iam_role_policy_attachment" {
  role       = aws_iam_role.dataloader_naptan_stopsinarea_v1_iam_role.name
  policy_arn = data.aws_iam_policy.AWSLambdaVPCAccessExecutionRole.arn
}

resource "aws_lambda_permission" "dataloader_naptan_stopsinarea_v1_cloudwatch_permission" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.dataloader_naptan_stopsinarea_v1_lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_day_at_0330.arn
}

resource "aws_lambda_permission" "dataloader_naptan_stopsinarea_v1_sns_permission" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.dataloader_naptan_stopsinarea_v1_lambda_function.function_name
  principal     = "sns.amazonaws.com"
  source_arn    = aws_sns_topic.production_elasticache_data_missing.arn
}

resource "aws_cloudwatch_event_target" "dataloader_naptan_stopsinarea_v1_invoke_every_day_at_0330" {
  depends_on = [
    aws_cloudwatch_event_rule.every_day_at_0330
  ]

  rule      = aws_cloudwatch_event_rule.every_day_at_0330.name
  target_id = "dataloader_naptan_stopsinarea_v1_lambda"
  arn       = aws_lambda_function.dataloader_naptan_stopsinarea_v1_lambda_function.arn
}

resource "aws_cloudwatch_log_group" "dataloader_naptan_stopsinarea_v1_cloudwatch_log_group" {
  name              = "/aws/lambda/dataloader_naptan_stopsinarea_v1"
  retention_in_days = var.cloudwatch_logs_retention_in_days
}

resource "aws_cloudwatch_metric_alarm" "dataloader_naptan_stopsinarea_v1_alarm_on_error_and_insufficient_data" {
  alarm_name          = "dataloader-naptan-stopsinarea-v1-alarm-if-error"
  alarm_description   = "Error occurred loading NaPTAN StopsInArea data"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  period              = 86400
  namespace           = "AWS/Lambda"
  metric_name         = "Errors"
  dimensions = {
    FunctionName = aws_lambda_function.dataloader_naptan_stopsinarea_v1_lambda_function.function_name
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
