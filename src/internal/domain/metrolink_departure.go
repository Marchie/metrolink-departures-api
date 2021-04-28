package domain

import (
	"time"
)

type MetrolinkDepartures struct {
	Departures  []*MetrolinkDeparture
	LastUpdated time.Time
}

type MetrolinkDeparture struct {
	AtcoCode    string
	Order       int
	Destination string
	Carriages   string
	Status      string
	Wait        string
	Platform    *string
	LastUpdated time.Time
}
