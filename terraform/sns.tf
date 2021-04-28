resource "aws_sns_topic" "production_notifications" {
  name = "production_notifications"
}

resource "aws_sns_topic_subscription" "production_notifications_subscription" {
  endpoint  = var.admin_email
  protocol  = "email"
  topic_arn = aws_sns_topic.production_notifications.arn
}

resource "aws_sns_topic" "production_elasticache_data_missing" {
  name = "production_elasticache_data_missing"
}

resource "aws_sns_topic_subscription" "production_elasticache_data_missing_subscription" {
  endpoint  = aws_lambda_function.dataloader_naptan_stopsinarea_v1_lambda_function.arn
  protocol  = "lambda"
  topic_arn = aws_sns_topic.production_elasticache_data_missing.arn
}
