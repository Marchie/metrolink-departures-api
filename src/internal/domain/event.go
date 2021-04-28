package domain

import "time"

type Event struct {
	StartTime time.Time
	Payload   string
}
