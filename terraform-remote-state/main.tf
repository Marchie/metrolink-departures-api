terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

# Create S3 bucket with random name
resource "aws_s3_bucket" "terraform_state" {
  bucket = var.tf_remote_state_bucket_name
  acl = "private"

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Apply policy
resource "aws_s3_bucket_policy" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  policy = jsonencode({
    Version = "2012-10-17",
    Id      = "RequireEncryption",
    Statement = [
      {
        Sid    = "RequireEncryptedTransport",
        Effect = "Deny",
        Action = [
          "s3:*"
        ],
        Resource = [
          "arn:aws:s3:::${aws_s3_bucket.terraform_state.bucket}/*"
        ],
        Condition = {
          Bool = {
            "aws:SecureTransport" : "false"
          }
        },
        Principal = "*"
      },
      {
        Sid    = "RequireEncryptedStorage",
        Effect = "Deny",
        Action = [
          "s3:PutObject"
        ],
        Resource = [
          "arn:aws:s3:::${aws_s3_bucket.terraform_state.bucket}/*"
        ],
        Condition = {
          StringNotEquals = {
            "s3:x-amz-server-side-encryption" : "AES256"
          }
        },
        Principal = "*"
      }
    ]
  })
}

# Create DynamoDB lock table
resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = var.tf_remote_state_table_name
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
