data "aws_iam_policy_document" "assume_apigateway_role_iam_policy_document" {
  version = "2012-10-17"

  statement {
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      identifiers = [
        "apigateway.amazonaws.com"
      ]
      type = "Service"
    }

    effect = "Allow"
  }
}

data "aws_iam_policy_document" "api_gateway_logging_iam_policy_document" {
  version = "2012-10-17"

  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
    ]

    resources = [
      "*"
    ]

    effect = "Allow"
  }
}

# Allow API Gateway to log
resource "aws_api_gateway_account" "tfgm_com" {
  cloudwatch_role_arn = aws_iam_role.api_gateway_cloudwatch.arn
}

resource "aws_iam_role" "api_gateway_cloudwatch" {
  name               = "api_gateway_cloudwatch_global"
  assume_role_policy = data.aws_iam_policy_document.assume_apigateway_role_iam_policy_document.json
}

resource "aws_iam_role_policy" "api_gateway_cloudwatch" {
  name   = "api_gateway_cloudwatch_global"
  policy = data.aws_iam_policy_document.api_gateway_logging_iam_policy_document.json
  role   = aws_iam_role.api_gateway_cloudwatch.id
}

# Root
resource "aws_api_gateway_rest_api" "tfgm_com" {
  name        = "TfGM.com API"
  description = "APIs to provide data for TfGM.com"
}

resource "aws_api_gateway_deployment" "tfgm_com_api_gateway_deployment" {
  rest_api_id = aws_api_gateway_integration.api_departures_metrolink_v1_lambda_api_gateway_integration.rest_api_id
}

resource "aws_api_gateway_stage" "tfgm_com_api_gateway_stage" {
  deployment_id = aws_api_gateway_deployment.tfgm_com_api_gateway_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.tfgm_com.id
  stage_name    = "tfgm_com"

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_departures_metrolink_v1_api_gateway_access_logs.arn
    format          = var.api_gateway_log_format_json
  }
}

resource "aws_api_gateway_method" "root_method" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_rest_api.tfgm_com.root_resource_id
  rest_api_id   = aws_api_gateway_rest_api.tfgm_com.id
}

resource "aws_api_gateway_integration" "root_integration" {
  rest_api_id = aws_api_gateway_method.root_method.rest_api_id
  resource_id = aws_api_gateway_method.root_method.resource_id
  http_method = aws_api_gateway_method.root_method.http_method
  type        = "MOCK"
}

resource "aws_api_gateway_integration_response" "root_integration_response" {
  http_method = aws_api_gateway_integration.root_integration.http_method
  resource_id = aws_api_gateway_integration.root_integration.resource_id
  rest_api_id = aws_api_gateway_integration.root_integration.rest_api_id
  status_code = 204
}

resource "aws_api_gateway_method_response" "root_response" {
  http_method = aws_api_gateway_integration_response.root_integration_response.http_method
  resource_id = aws_api_gateway_integration_response.root_integration_response.resource_id
  rest_api_id = aws_api_gateway_integration_response.root_integration_response.rest_api_id
  status_code = aws_api_gateway_integration_response.root_integration_response.status_code
}

# /departures
resource "aws_api_gateway_resource" "departures_resource" {
  rest_api_id = aws_api_gateway_rest_api.tfgm_com.id
  parent_id   = aws_api_gateway_rest_api.tfgm_com.root_resource_id
  path_part   = "departures"
}

resource "aws_api_gateway_method" "departures_method" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.departures_resource.id
  rest_api_id   = aws_api_gateway_resource.departures_resource.rest_api_id
}

resource "aws_api_gateway_integration" "departures_integration" {
  http_method = aws_api_gateway_method.departures_method.http_method
  resource_id = aws_api_gateway_method.departures_method.resource_id
  rest_api_id = aws_api_gateway_method.departures_method.rest_api_id
  type        = "MOCK"
}

resource "aws_api_gateway_integration_response" "departures_integration_response" {
  http_method = aws_api_gateway_integration.departures_integration.http_method
  resource_id = aws_api_gateway_integration.departures_integration.resource_id
  rest_api_id = aws_api_gateway_integration.departures_integration.rest_api_id
  status_code = 204
}

resource "aws_api_gateway_method_response" "departures_method_response" {
  http_method = aws_api_gateway_integration_response.departures_integration_response.http_method
  resource_id = aws_api_gateway_integration_response.departures_integration_response.resource_id
  rest_api_id = aws_api_gateway_integration_response.departures_integration_response.rest_api_id
  status_code = aws_api_gateway_integration_response.departures_integration_response.status_code
}

# /departures/metrolink

resource "aws_api_gateway_resource" "departures_metrolink_resource" {
  parent_id   = aws_api_gateway_resource.departures_resource.id
  path_part   = "metrolink"
  rest_api_id = aws_api_gateway_resource.departures_resource.rest_api_id
}

resource "aws_api_gateway_method" "departures_metrolink_method" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.departures_metrolink_resource.id
  rest_api_id   = aws_api_gateway_resource.departures_metrolink_resource.rest_api_id
}

resource "aws_api_gateway_integration" "departures_metrolink_integration" {
  http_method = aws_api_gateway_method.departures_metrolink_method.http_method
  resource_id = aws_api_gateway_method.departures_metrolink_method.resource_id
  rest_api_id = aws_api_gateway_method.departures_metrolink_method.rest_api_id
  type        = "MOCK"
}

resource "aws_api_gateway_integration_response" "departures_metrolink_integration_response" {
  http_method = aws_api_gateway_integration.departures_metrolink_integration.http_method
  resource_id = aws_api_gateway_integration.departures_metrolink_integration.resource_id
  rest_api_id = aws_api_gateway_integration.departures_metrolink_integration.rest_api_id
  status_code = 204
}

resource "aws_api_gateway_method_response" "departures_metrolink_method_response" {
  http_method = aws_api_gateway_integration_response.departures_metrolink_integration_response.http_method
  resource_id = aws_api_gateway_integration_response.departures_metrolink_integration_response.resource_id
  rest_api_id = aws_api_gateway_integration_response.departures_metrolink_integration_response.rest_api_id
  status_code = aws_api_gateway_integration_response.departures_metrolink_integration_response.status_code
}

# /departures/metrolink/v1

resource "aws_api_gateway_resource" "departures_metrolink_v1_resource" {
  parent_id   = aws_api_gateway_resource.departures_metrolink_resource.id
  path_part   = "v1"
  rest_api_id = aws_api_gateway_resource.departures_metrolink_resource.rest_api_id
}

resource "aws_api_gateway_method" "departures_metrolink_v1_method" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.departures_metrolink_v1_resource.id
  rest_api_id   = aws_api_gateway_resource.departures_metrolink_v1_resource.rest_api_id
}

resource "aws_api_gateway_integration" "departures_metrolink_v1_integration" {
  http_method = aws_api_gateway_method.departures_metrolink_v1_method.http_method
  resource_id = aws_api_gateway_method.departures_metrolink_v1_method.resource_id
  rest_api_id = aws_api_gateway_method.departures_metrolink_v1_method.rest_api_id
  type        = "MOCK"
}

resource "aws_api_gateway_integration_response" "departures_metrolink_v1_integration_response" {
  http_method = aws_api_gateway_integration.departures_metrolink_v1_integration.http_method
  resource_id = aws_api_gateway_integration.departures_metrolink_v1_integration.resource_id
  rest_api_id = aws_api_gateway_integration.departures_metrolink_v1_integration.rest_api_id
  status_code = 204
}

resource "aws_api_gateway_method_response" "departures_metrolink_v1_method_response" {
  http_method = aws_api_gateway_integration_response.departures_metrolink_v1_integration_response.http_method
  resource_id = aws_api_gateway_integration_response.departures_metrolink_v1_integration_response.resource_id
  rest_api_id = aws_api_gateway_integration_response.departures_metrolink_v1_integration_response.rest_api_id
  status_code = aws_api_gateway_integration_response.departures_metrolink_v1_integration_response.status_code
}
