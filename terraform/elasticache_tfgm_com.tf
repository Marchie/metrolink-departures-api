resource "aws_elasticache_cluster" "tfgm_com" {
  cluster_id           = "production-tfgm-com"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis6.x"
  engine_version       = "6.0.5"
  port                 = 6379
  maintenance_window   = "mon:01:00-mon:03:00"

  apply_immediately = true

  security_group_ids = [
    aws_security_group.production_elasticache_sg.id
  ]

  subnet_group_name = aws_elasticache_subnet_group.production_persistence.name

  notification_topic_arn = aws_sns_topic.production_notifications.arn

  tags = {
    Environment = "Production"
  }
}

resource "aws_elasticache_subnet_group" "production_persistence" {
  name = "production-persistence"
  subnet_ids = [
    aws_subnet.production_persistence.id
  ]
}

resource "aws_cloudwatch_metric_alarm" "production_elasticache_missing-data" {
  alarm_name          = "production-elasticache-missing-data"
  alarm_description   = "ElastiCache has lower than expected number of records"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 1
  period              = 300
  namespace           = "AWS/ElastiCache"
  metric_name         = "CurrItems"
  dimensions = {
    CacheClusterId = aws_elasticache_cluster.tfgm_com.cluster_id
    CacheNodeId    = aws_elasticache_cluster.tfgm_com.cache_nodes.0.id
  }
  statistic       = "Average"
  threshold       = 89000
  actions_enabled = true
  alarm_actions = [
    aws_sns_topic.production_notifications.arn,
    aws_sns_topic.production_elasticache_data_missing.arn
  ]
  ok_actions = [
    aws_sns_topic.production_notifications.arn
  ]
  insufficient_data_actions = [
    aws_sns_topic.production_notifications.arn
  ]
}
