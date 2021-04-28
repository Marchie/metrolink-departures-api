# Metrolink Departures API

A prototype application which sources departure data from the TfGM Metrolinks API and presents it in a more consumable
way.

## Architecture

See [metrolink-departures-api.drawio](metrolink-departures-api.drawio) for an architecture diagram.

The application consists of four Lambda functions:

* [api-departures-metrolink-v1](src/cmd/api/departures/metrolink/v1/README.md)
* [dataloader-departures-metrolink-v1](src/cmd/dataloader/departures/metrolink/v1/README.md)
* [dataloader-naptan-stopsinarea-v1](src/cmd/dataloader/naptan/stopsinarea/v1/README.md)
* [scheduler-departures-metrolink-v1](src/cmd/scheduler/departures/metrolink/v1/README.md)
