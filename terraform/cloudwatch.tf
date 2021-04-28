resource "aws_cloudwatch_event_rule" "every_one_minute" {
  name                = "every-one-minute"
  description         = "Fires once per minute"
  schedule_expression = "cron(* * * * ? *)"
}

resource "aws_cloudwatch_event_rule" "every_day_at_0330" {
  name                = "every-day-at-0330"
  description         = "Fires once per day at 03:30"
  schedule_expression = "cron(30 3 * * ? *)"
}
