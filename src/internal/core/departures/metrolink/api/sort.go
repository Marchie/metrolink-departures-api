package api

import (
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"strconv"
)

type ByWaitStatusDestinationOrderPlatformCarriages []*domain.MetrolinkDeparture

func (s ByWaitStatusDestinationOrderPlatformCarriages) Len() int {
	return len(s)
}

func (s ByWaitStatusDestinationOrderPlatformCarriages) Less(i, j int) bool {
	if s[i].Wait == s[j].Wait {
		if s[i].Status == s[j].Status {
			if s[i].Order == s[j].Order {
				if s[i].Destination == s[j].Destination {
					if s[i].Platform != nil && s[j].Platform != nil && *s[i].Platform != *s[j].Platform {
						return *s[i].Platform < *s[j].Platform
					}

					return s[i].Carriages > s[j].Carriages
				}

				return s[i].Destination < s[j].Destination
			}

			return s[i].Order < s[j].Order
		}

		// The possible statuses are:
		// Departing - the tram is setting off from the stop
		// Arrived - the tram is at the stop
		// Due - the tram is en route to the stop

		if s[i].Status == "Departing" {
			return true
		}

		if s[j].Status == "Departing" {
			return false
		}

		if s[i].Status == "Arrived" {
			return true
		}

		return false
	}

	// Wait value is potentially "DELAY"
	// Display "DELAY" values last
	iWait, err := strconv.Atoi(s[i].Wait)
	if err != nil {
		return false
	}

	jWait, err := strconv.Atoi(s[j].Wait)
	if err != nil {
		return true
	}

	return iWait < jWait
}

func (s ByWaitStatusDestinationOrderPlatformCarriages) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
