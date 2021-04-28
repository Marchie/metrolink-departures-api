variable "aws_region" {
  description = "The AWS region for the remote Terraform state"
  type        = string
}

variable "aws_profile" {
  description = "The AWS profile to use"
  type        = string
}

variable "tf_remote_state_bucket_name" {
  description = "The S3 bucket name for the remote Terraform state"
  type        = string
}

variable "tf_remote_state_table_name" {
  description = "The DynamoDB table name for the remote Terraform state"
  type        = string
}
