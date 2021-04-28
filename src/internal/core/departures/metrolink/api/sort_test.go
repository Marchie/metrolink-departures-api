package api_test

import (
	"github.com/Marchie/tf-experiment/lambda/internal/core/departures/metrolink/api"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func givenASliceOfMetrolinkDeparture(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	return []*domain.MetrolinkDeparture{
		{
			Order:       0,
			Destination: "Altrincham",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "5",
			Platform:    nil,
		},
		{
			Order:       1,
			Destination: "East Didsbury",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "7",
			Platform:    nil,
		},
	}
}

func TestByWaitStatusDestinationOrderPlatformCarriages_Len(t *testing.T) {
	t.Run(`Given a slice of two domain.MetrolinkDeparture
When Len() is called
Then the length of the slice is returned`, func(t *testing.T) {
		// Given
		metrolinkDepartures := givenASliceOfMetrolinkDeparture(t)

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Len()

		// Then
		assert.Equal(t, result, 2)
	})
}

func TestByWaitStatusDestinationOrderPlatformCarriages_Swap(t *testing.T) {
	t.Run(`Given a slice of two domain.MetrolinkDeparture
When Swap() is called with the two indices
Then the position of these items in the slice is swapped`, func(t *testing.T) {
		// Given
		metrolinkDepartures := givenASliceOfMetrolinkDeparture(t)

		// When
		api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Swap(0, 1)

		// Then
		assert.Equal(t, metrolinkDepartures[0].Destination, "East Didsbury")
		assert.Equal(t, metrolinkDepartures[1].Destination, "Altrincham")
	})
}

func TestByWaitStatusDestinationOrderPlatformCarriages_Less(t *testing.T) {
	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has a lower wait time value than the second departure
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "2",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has a larger wait time value than the second departure
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "2",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has a wait time of DELAY
And the second departure has a numerical wait time
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "DELAY",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has a numerical wait time
And the second departure has a wait time of DELAY
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "DELAY",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has a status of "Departing"
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Departing",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has a status of "Arrived"
And the second departure has a status of "Departing"
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Arrived",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Departing",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has a status of "Arrived"
And the second departure has a status of "Due"
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Arrived",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has a status of "Due"
And the second departure has a status of "Arrived"
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Single",
				Status:      "Arrived",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has a lower order value than the second departure
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       1,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has a higher order value than the second departure
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       1,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "East Didsbury",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal order to the second departure
And the first departure has a destination which is alphabetically before the second departure destination
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Ashton",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal order to the second departure
And the first departure has a destination which is alphabetically after the second departure destination
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Ashton",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal destination to the second departure
And the first departure has an equal order value to the second departure
And the first departure has an equal platform to the second departure
And the first departure has is a Single carriage and the second departure is a Double carriage
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal destination to the second departure
And the first departure has an equal order value to the second departure
And the first departure has an equal platform to the second departure
And the first departure has is a Double carriage and the second departure is a Single carriage
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    nil,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal destination to the second departure
And the first departure has an equal order value to the second departure
And the first departure platform is alphabetically before the second departure
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		platformA := "A"
		platformB := "B"

		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platformA,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platformB,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal destination to the second departure
And the first departure has an equal order value to the second departure
And the first departure platform is alphabetically after the second departure
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		platformA := "A"
		platformB := "B"

		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platformB,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platformA,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal destination to the second departure
And the first departure has an equal order value to the second departure
And the first departure platform is numerically before the second departure
When Less() is called with the two indices
Then Less() returns true`, func(t *testing.T) {
		// Given
		platform1 := "1"
		platform2 := "2"

		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platform1,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platform2,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, true, result)
	})

	t.Run(`Given a slice of domain.MetrolinkDeparture
And the first departure has an equal wait time value to the second departure
And the first departure has an equal status to the second departure
And the first departure has an equal destination to the second departure
And the first departure has an equal order value to the second departure
And the first departure platform is numerically before the second departure
When Less() is called with the two indices
Then Less() returns false`, func(t *testing.T) {
		// Given
		platform1 := "1"
		platform2 := "2"

		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platform2,
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "1",
				Platform:    &platform1,
			},
		}

		// When
		result := api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures).Less(0, 1)

		// Then
		assert.Equal(t, false, result)
	})
}

func TestByWaitStatusDestinationOrderPlatformCarriages_Sorting(t *testing.T) {
	t.Run(`Given a slice of MetrolinkDeparture
When sorted ByWaitStatusDestinationOrderPlatformCarriages
Then the MetrolinkDeparture slice is sorted correctly`, func(t *testing.T) {
		// Given
		metrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "2",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       1,
				Destination: "Eccles via MediaCityUK",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "9",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       2,
				Destination: "Manchester Airport",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "17",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       0,
				Destination: "Bury",
				Carriages:   "Double",
				Status:      "Arrived",
				Wait:        "0",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       1,
				Destination: "Rochdale via Oldham",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "5",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       2,
				Destination: "Victoria",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "9",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
		}

		// When
		sort.Sort(api.ByWaitStatusDestinationOrderPlatformCarriages(metrolinkDepartures))

		// Then
		expectedMetrolinkDepartures := []*domain.MetrolinkDeparture{
			{
				Order:       0,
				Destination: "Bury",
				Carriages:   "Double",
				Status:      "Arrived",
				Wait:        "0",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       0,
				Destination: "Altrincham",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "2",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       1,
				Destination: "Rochdale via Oldham",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "5",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       1,
				Destination: "Eccles via MediaCityUK",
				Carriages:   "Double",
				Status:      "Due",
				Wait:        "9",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       2,
				Destination: "Victoria",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "9",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
			{
				Order:       2,
				Destination: "Manchester Airport",
				Carriages:   "Single",
				Status:      "Due",
				Wait:        "17",
				Platform:    nil,
				LastUpdated: time.Time{},
			},
		}

		assert.EqualValues(t, expectedMetrolinkDepartures, metrolinkDepartures)
	})
}
