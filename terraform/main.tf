terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  backend "s3" {
    bucket         = "marchie-terraform-remote-state"
    key            = "demo/terraform.tfstate"
    region         = "eu-west-1"
    profile        = "marchie_root"
    dynamodb_table = "marchie-terraform-remote-lock"
    encrypt        = true
  }
}

provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
}

resource "aws_s3_bucket" "build_artifacts" {
  bucket = var.build_artifacts_bucket_name
  acl    = "private"
}
