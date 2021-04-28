# dataloader-naptan-stopsinarea-v1

A Lambda function which retrieves NaPTAN data in CSV format and stores the relationship between a `StopAreaCode` and
`AtcoCode` in an ElastiCache repository. The data stored in the repository is used as source data for the 
[api-departures-metrolink-v1 Lambda function](../../../../api/departures/metrolink/v1/README.md).

