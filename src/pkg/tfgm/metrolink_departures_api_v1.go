package tfgm

import "time"

type MetrolinkDepartures struct {
	RequestedLocation string                `json:"requestedLocation"`
	Departures        []*MetrolinkDeparture `json:"departures"`
	LastUpdated       time.Time             `json:"lastUpdated"`
}

type MetrolinkDeparture struct {
	AtcoCode    string    `json:"atcoCode"`
	Sequence    int       `json:"sequence"`
	Destination string    `json:"destination"`
	Status      string    `json:"status"`
	Wait        string    `json:"wait"`
	Carriages   string    `json:"carriages"`
	Platform    *string   `json:"platform,omitempty"`
	LastUpdated time.Time `json:"lastUpdated"`
}
