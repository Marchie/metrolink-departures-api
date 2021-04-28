variable "admin_email" {
  description = "Email address to send notifications and alerts to"
  type        = string
}

variable "aws_region" {
  description = "The AWS region for the remote Terraform state"
  type        = string
}

variable "aws_profile" {
  description = "The AWS profile to use"
  type        = string
}

variable "build_artifacts_bucket_name" {
  description = "The S3 bucket name for build artifacts"
  type        = string
}

variable "cloudwatch_logs_retention_in_days" {
  description = "The number of days to retain CloudWatch logs for"
  type        = number
  default     = 7
}

variable "tf_remote_state_bucket_name" {
  description = "The S3 bucket name for the remote Terraform state"
  type        = string
}

variable "tf_remote_state_key" {
  description = "The S3 key for the remote Terraform state"
  type        = string
}

variable "tf_remote_state_table_name" {
  description = "The DynamoDB table name for the remote Terraform state"
  type        = string
}

variable "departures_metrolink_v1_tfgm_developer_metrolinks_url" {
  description = "The URL for the TfGM Developer API Metrolinks endpoint"
  type        = string
  default     = "https://api.tfgm.com/tfgm/Metrolinks"
}

variable "departures_metrolink_v1_tfgm_developer_api_key" {
  description = "The API Key for the TfGM Developer API"
  type        = string
}

variable "api_gateway_log_format_json" {
  description = "Log format for API Gateway"
  type        = string
  default     = "{\"requestId\":\"$context.requestId\", \"ip\": \"$context.identity.sourceIp\", \"caller\":\"$context.identity.caller\", \"user\":\"$context.identity.user\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\", \"resourcePath\":\"$context.resourcePath\", \"status\":\"$context.status\", \"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\"}"
}

variable "naptan_csv_zip_url" {
  description = "The URL for the NaPTAN CSV Zip archive"
  type        = string
  default     = "https://naptan.app.dft.gov.uk/DataRequest/Naptan.ashx?format=csv"
}
