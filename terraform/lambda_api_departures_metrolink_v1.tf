# /departures/metrolink/v1/{stopAreaCodeOrAtcoCode}

resource "aws_api_gateway_resource" "api_departures_metrolink_v1_lambda_api_gateway_resource" {
  rest_api_id = aws_api_gateway_resource.departures_metrolink_v1_resource.rest_api_id
  parent_id   = aws_api_gateway_resource.departures_metrolink_v1_resource.id
  path_part   = "{stopAreaCodeOrAtcoCode}"
}

resource "aws_api_gateway_method" "api_departures_metrolink_v1_lambda_api_gateway_method" {
  authorization = "NONE"
  http_method   = "GET"
  resource_id   = aws_api_gateway_resource.api_departures_metrolink_v1_lambda_api_gateway_resource.id
  rest_api_id   = aws_api_gateway_resource.api_departures_metrolink_v1_lambda_api_gateway_resource.rest_api_id
}

resource "aws_api_gateway_integration" "api_departures_metrolink_v1_lambda_api_gateway_integration" {
  http_method             = aws_api_gateway_method.api_departures_metrolink_v1_lambda_api_gateway_method.http_method
  resource_id             = aws_api_gateway_method.api_departures_metrolink_v1_lambda_api_gateway_method.resource_id
  rest_api_id             = aws_api_gateway_method.api_departures_metrolink_v1_lambda_api_gateway_method.rest_api_id
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = aws_lambda_function.api_departures_metrolink_v1_lambda_function.invoke_arn
  passthrough_behavior    = "WHEN_NO_MATCH"
}


resource "aws_lambda_function" "api_departures_metrolink_v1_lambda_function" {
  depends_on = [
    aws_iam_role_policy_attachment.api_departures_metrolink_v1_iam_role_policy_attachment,
    aws_cloudwatch_log_group.api_departures_metrolink_v1_cloudwatch_log_group,
    aws_elasticache_cluster.tfgm_com
  ]

  function_name = "api_departures_metrolink_v1"

  s3_bucket = "marchie-build-artifacts"
  s3_key    = "v0.0.1/api_departures_metrolink_v1.zip"

  handler = "api_departures_metrolink_v1"
  runtime = "go1.x"

  timeout     = 1
  memory_size = 320

  reserved_concurrent_executions = 5

  role = aws_iam_role.api_departures_metrolink_v1_iam_role.arn

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
      REDIS_METROLINK_DEPARTURES_SERVER_ADDRESS                = "${aws_elasticache_cluster.tfgm_com.cache_nodes.0.address}:${aws_elasticache_cluster.tfgm_com.cache_nodes.0.port}"
      REDIS_METROLINK_DEPARTURES_SERVICE_STATUS_SERVER_ADDRESS = "${aws_elasticache_cluster.tfgm_com.cache_nodes.0.address}:${aws_elasticache_cluster.tfgm_com.cache_nodes.0.port}"
      REDIS_STOPS_IN_AREA_SERVER_ADDRESS                       = "${aws_elasticache_cluster.tfgm_com.cache_nodes.0.address}:${aws_elasticache_cluster.tfgm_com.cache_nodes.0.port}"
      STOP_AREA_CODE_OR_ATCO_CODE_API_GATEWAY_PATH_PARAMETER   = "stopAreaCodeOrAtcoCode"
    }
  }
}

resource "aws_iam_role" "api_departures_metrolink_v1_iam_role" {
  name = "api_departures_metrolink_v1"

  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role_iam_policy_document.json
}

resource "aws_iam_role_policy_attachment" "api_departures_metrolink_v1_iam_role_policy_attachment" {
  role = aws_iam_role.api_departures_metrolink_v1_iam_role.name

  policy_arn = data.aws_iam_policy.AWSLambdaVPCAccessExecutionRole.arn
}

resource "aws_lambda_permission" "api_departures_metrolink_v1_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api_departures_metrolink_v1_lambda_function.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.tfgm_com.execution_arn}/*/*"
}

resource "aws_api_gateway_method_settings" "api_departures_metrolink_v1_api_gateway_method_settings" {
  method_path = "*/*"
  rest_api_id = aws_api_gateway_method.api_departures_metrolink_v1_lambda_api_gateway_method.rest_api_id
  stage_name  = aws_api_gateway_stage.tfgm_com_api_gateway_stage.stage_name

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }
}

resource "aws_cloudwatch_log_group" "api_departures_metrolink_v1_api_gateway_access_logs" {
  name              = "API-Gateway-Access-Logs_${aws_api_gateway_rest_api.tfgm_com.id}/api_departures_metrolink_v1"
  retention_in_days = var.cloudwatch_logs_retention_in_days
}

resource "aws_cloudwatch_log_group" "api_departures_metrolink_v1_api_gateway_execution_logs" {
  name              = "API-Gateway-Execution-Logs_${aws_api_gateway_rest_api.tfgm_com.id}/tfgm_com"
  retention_in_days = var.cloudwatch_logs_retention_in_days
}

resource "aws_cloudwatch_log_group" "api_departures_metrolink_v1_cloudwatch_log_group" {
  name              = "/aws/lambda/api_departures_metrolink_v1"
  retention_in_days = var.cloudwatch_logs_retention_in_days
}

resource "aws_cloudwatch_metric_alarm" "api_departures_metrolink_v1_alarm_on_error" {
  alarm_name          = "api-departures-metrolink-v1-alarm-if-error"
  alarm_description   = "Error occurred with Metrolink departures API"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = 1
  period              = 300
  namespace           = "AWS/Lambda"
  metric_name         = "Errors"
  dimensions = {
    FunctionName = aws_lambda_function.api_departures_metrolink_v1_lambda_function.function_name
  }
  statistic       = "Sum"
  threshold       = 1
  actions_enabled = true
  alarm_actions = [
    aws_sns_topic.production_notifications.arn
  ]
}

output "api_departures_metrolink_v1_url" {
  value = "${aws_api_gateway_deployment.tfgm_com_api_gateway_deployment.invoke_url}${aws_api_gateway_stage.tfgm_com_api_gateway_stage.stage_name}"
}
