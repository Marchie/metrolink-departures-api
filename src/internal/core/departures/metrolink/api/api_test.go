package api

import (
	"context"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	mock_repository "github.com/Marchie/tf-experiment/lambda/internal/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func mockLogger(t *testing.T) *zap.Logger {
	t.Helper()

	zapCore, _ := observer.New(zapcore.ErrorLevel)
	return zap.New(zapCore)
}

func readJson(t *testing.T, rc io.ReadCloser) string {
	t.Helper()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	return string(b)
}

func givenMetrolinkDeparturesForAtcoCode9400ZZMASTP1(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	platform := "D"

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       0,
			Destination: "Rochdale",
			Carriages:   "Single",
			Status:      "Departing",
			Wait:        "0",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       1,
			Destination: "Victoria",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "4",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       2,
			Destination: "Rochdale",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "12",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
	}
}

func givenMetrolinkDeparturesForAtcoCode9400ZZMASTP2(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	platform := "C"

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       0,
			Destination: "Piccadilly",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "1",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       1,
			Destination: "Ashton-under-Lyne",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "5",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       2,
			Destination: "Bury",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "9",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
	}
}

func givenMetrolinkDeparturesForAtcoCode9400ZZMASTP3(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	platform := "B"

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP3",
			Order:       0,
			Destination: "East Didsbury",
			Carriages:   "Double",
			Status:      "Arrived",
			Wait:        "0",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP3",
			Order:       1,
			Destination: "East Didsbury",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "12",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP3",
			Order:       2,
			Destination: "East Didsbury",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "24",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
	}
}

func givenMetrolinkDeparturesForAtcoCode9400ZZMASTP4(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	platform := "A"

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP4",
			Order:       0,
			Destination: "Altrincham",
			Carriages:   "Double",
			Status:      "Departing",
			Wait:        "0",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP4",
			Order:       1,
			Destination: "Eccles via MediaCityUK",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "3",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP4",
			Order:       2,
			Destination: "Manchester Airport",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "7",
			Platform:    &platform,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
	}
}

func givenMetrolinkDeparturesForAtcoCode9400ZZMAMKT1(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMAMKT1",
			Order:       0,
			Destination: "Rochdale",
			Carriages:   "Single",
			Status:      "Departing",
			Wait:        "0",
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMAMKT1",
			Order:       1,
			Destination: "Victoria",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "4",
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMAMKT1",
			Order:       2,
			Destination: "Rochdale",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "12",
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
	}
}

func thenExpectDeparturesFor940GZZMASTP(t *testing.T) []*domain.MetrolinkDeparture {
	t.Helper()

	platformA := "A"
	platformB := "B"
	platformC := "C"
	platformD := "D"

	return []*domain.MetrolinkDeparture{
		{
			AtcoCode:    "9400ZZMASTP4",
			Order:       0,
			Destination: "Altrincham",
			Carriages:   "Double",
			Status:      "Departing",
			Wait:        "0",
			Platform:    &platformA,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       0,
			Destination: "Rochdale",
			Carriages:   "Single",
			Status:      "Departing",
			Wait:        "0",
			Platform:    &platformD,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP3",
			Order:       0,
			Destination: "East Didsbury",
			Carriages:   "Double",
			Status:      "Arrived",
			Wait:        "0",
			Platform:    &platformB,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       0,
			Destination: "Piccadilly",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "1",
			Platform:    &platformC,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP4",
			Order:       1,
			Destination: "Eccles via MediaCityUK",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "3",
			Platform:    &platformA,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       1,
			Destination: "Victoria",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "4",
			Platform:    &platformD,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       1,
			Destination: "Ashton-under-Lyne",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "5",
			Platform:    &platformC,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP4",
			Order:       2,
			Destination: "Manchester Airport",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "7",
			Platform:    &platformA,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP2",
			Order:       2,
			Destination: "Bury",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "9",
			Platform:    &platformC,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP3",
			Order:       1,
			Destination: "East Didsbury",
			Carriages:   "Single",
			Status:      "Due",
			Wait:        "12",
			Platform:    &platformB,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP1",
			Order:       2,
			Destination: "Rochdale",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "12",
			Platform:    &platformD,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
		{
			AtcoCode:    "9400ZZMASTP3",
			Order:       2,
			Destination: "East Didsbury",
			Carriages:   "Double",
			Status:      "Due",
			Wait:        "24",
			Platform:    &platformB,
			LastUpdated: time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC),
		},
	}
}

func givenCurrentTimeFunc(t *testing.T) func() time.Time {
	t.Helper()

	return func() time.Time {
		return time.Date(2021, time.April, 6, 21, 37, 30, 0, time.UTC)
	}
}

func givenLastUpdatedTime(t *testing.T) *time.Time {
	t.Helper()

	lastUpdatedTime := time.Date(2021, time.April, 6, 21, 37, 19, 0, time.UTC)

	return &lastUpdatedTime
}

func givenStaleLastUpdatedTime(t *testing.T) *time.Time {
	t.Helper()

	staleLastUpdatedTime := time.Date(2021, time.April, 6, 21, 36, 44, 0, time.UTC)

	return &staleLastUpdatedTime
}

func givenStaleDataThreshold(t *testing.T) time.Duration {
	t.Helper()

	return time.Second * 45
}

func givenTimeLocation(t *testing.T) *time.Location {
	t.Helper()

	loc, _ := time.LoadLocation("Europe/London")

	return loc
}

func thenExpectJsonError(t *testing.T, errorMsg string, requestedLocation string) string {
	t.Helper()

	return "{\n\t\"error\": \"" + errorMsg + "\",\n\t\"requestedLocation\": \"" + requestedLocation + "\"\n}\n"
}

func thenExpectJsonDeparturesWithoutPlatform(t *testing.T, requestedLocation string, departures []*domain.MetrolinkDeparture, lastUpdated *time.Time) string {
	t.Helper()

	expString := "{\n\t\"requestedLocation\": \"" + requestedLocation + "\",\n\t\"departures\": ["
	for sequence, departure := range departures {
		if sequence > 0 {
			expString += ","
		}
		expString += fmt.Sprintf("\n\t\t{\n\t\t\t\"atcoCode\": \"%s\",\n\t\t\t\"sequence\": %d,\n\t\t\t\"destination\": \"%s\",\n\t\t\t\"status\": \"%s\",\n\t\t\t\"wait\": \"%s\",\n\t\t\t\"carriages\": \"%s\",\n\t\t\t\"lastUpdated\": \"%s\"\n\t\t}", departure.AtcoCode, sequence, departure.Destination, departure.Status, departure.Wait, departure.Carriages, departure.LastUpdated.In(givenTimeLocation(t)).Format(time.RFC3339))
	}
	expString += "\n\t],\n\t\"lastUpdated\": \"" + lastUpdated.In(givenTimeLocation(t)).Format(time.RFC3339) + "\"\n}\n"

	return expString
}

func thenExpectJsonDeparturesWithPlatform(t *testing.T, requestedLocation string, departures []*domain.MetrolinkDeparture, lastUpdated *time.Time) string {
	t.Helper()

	expString := "{\n\t\"requestedLocation\": \"" + requestedLocation + "\",\n\t\"departures\": ["
	for sequence, departure := range departures {
		if sequence > 0 {
			expString += ","
		}
		expString += fmt.Sprintf("\n\t\t{\n\t\t\t\"atcoCode\": \"%s\",\n\t\t\t\"sequence\": %d,\n\t\t\t\"destination\": \"%s\",\n\t\t\t\"status\": \"%s\",\n\t\t\t\"wait\": \"%s\",\n\t\t\t\"carriages\": \"%s\",\n\t\t\t\"platform\": \"%s\",\n\t\t\t\"lastUpdated\": \"%s\"\n\t\t}", departure.AtcoCode, sequence, departure.Destination, departure.Status, departure.Wait, departure.Carriages, *departure.Platform, departure.LastUpdated.In(givenTimeLocation(t)).Format(time.RFC3339))
	}
	expString += "\n\t],\n\t\"lastUpdated\": \"" + lastUpdated.In(givenTimeLocation(t)).Format(time.RFC3339) + "\"\n}\n"

	return expString
}

func thenExpectJsonWithEmptyDeparturesSlice(t *testing.T, requestedLocation string, lastUpdated *time.Time) string {
	t.Helper()

	return "{\n\t\"requestedLocation\": \"" + requestedLocation + "\",\n\t\"departures\": [],\n\t\"lastUpdated\": \"" + lastUpdated.In(givenTimeLocation(t)).Format(time.RFC3339) + "\"\n}\n"
}

func TestApi_Json(t *testing.T) {
	t.Run(`Given an invalid Metrolink StopAreaCode or AtcoCode is requested
When Json is called
Then an error JSON response is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)

		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)

		invalidMetrolinkStopAreaCodeOrAtcoCode := "1800SB18811"

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, invalidMetrolinkStopAreaCodeOrAtcoCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, thenExpectJsonError(t, "invalid StopAreaCode or AtcoCode", invalidMetrolinkStopAreaCodeOrAtcoCode), readJson(t, rc))
	})

	t.Run(`Given an error occurs fetching the system status
When Json is called
Then an error JSON response is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		validMetrolinkStopAreaCode := "940GZZMASTP"

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)

		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(nil, redis.ErrNil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkStopAreaCode)

		// Then
		assert.Nil(t, rc)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting Metrolink departures system status: redigo: nil returned")
	})

	t.Run(`Given the system status shows that the last updated time breaches the stale data threshold
When Json is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		validMetrolinkStopAreaCode := "940GZZMASTP"

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)

		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)

		staleLastUpdatedTime := givenStaleLastUpdatedTime(t)
		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(staleLastUpdatedTime, nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkStopAreaCode)

		// Then
		assert.NotNil(t, rc)
		assert.Equal(t, http.StatusBadGateway, statusCode)
		assert.Nil(t, err)
		assert.Equal(t, thenExpectJsonError(t, "Metrolink departures data is outdated: last updated at 2021-04-06T21:36:44Z", validMetrolinkStopAreaCode), readJson(t, rc))
	})

	t.Run(`Given an potentially valid StopAreaCode is requested
And the StopAreaCode contains no AtcoCodes
When Json is called
Then an error JSON response is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		invalidMetrolinkStopAreaCode := "940GZZMAXXX"

		wrappedRedisErrNil := errors.Wrap(redis.ErrNil, "FUBAR")
		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)
		stopsInAreaGetter.EXPECT().GetStopsInArea(ctx, invalidMetrolinkStopAreaCode).Return(nil, wrappedRedisErrNil)

		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(givenLastUpdatedTime(t), nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, invalidMetrolinkStopAreaCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, thenExpectJsonError(t, "invalid StopAreaCode", invalidMetrolinkStopAreaCode), readJson(t, rc))
	})

	t.Run(`Given a valid Metrolink AtcoCode is requested
When Json is called
Then departures are returned for that AtcoCode`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)

		validMetrolinkAtcoCode := "9400ZZMASTP1"

		metrolinkDeparturesForAtcoCode := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP1(t)
		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, validMetrolinkAtcoCode).Return(metrolinkDeparturesForAtcoCode, nil)

		lastUpdatedTime := givenLastUpdatedTime(t)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(lastUpdatedTime, nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkAtcoCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, thenExpectJsonDeparturesWithPlatform(t, validMetrolinkAtcoCode, metrolinkDeparturesForAtcoCode, lastUpdatedTime), readJson(t, rc))
	})

	t.Run(`Given a valid Metrolink AtcoCode is requested
And the AtcoCode does not have an associated platform
When Json is called
Then departures are returned for that AtcoCode
And there is no platform value for the departures`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)

		validMetrolinkAtcoCode := "9400ZZMASTP1"

		metrolinkDeparturesForAtcoCode := givenMetrolinkDeparturesForAtcoCode9400ZZMAMKT1(t)
		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, validMetrolinkAtcoCode).Return(metrolinkDeparturesForAtcoCode, nil)

		lastUpdatedTime := givenLastUpdatedTime(t)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(lastUpdatedTime, nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkAtcoCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, thenExpectJsonDeparturesWithoutPlatform(t, validMetrolinkAtcoCode, metrolinkDeparturesForAtcoCode, lastUpdatedTime), readJson(t, rc))
	})

	t.Run(`Given a valid Metrolink StopAreaCode is requested
When Json is called
Then sorted departures are returned for each AtcoCode in that StopAreaCode`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		validMetrolinkStopAreaCode := "940GZZMASTP"

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)
		stopsInAreaGetter.EXPECT().GetStopsInArea(ctx, validMetrolinkStopAreaCode).Return([]string{"9400ZZMASTP", "9400ZZMASTP1", "9400ZZMASTP2", "9400ZZMASTP3", "9400ZZMASTP4"}, nil)

		metrolinkDeparturesForAtcoCode9400ZZMASTP1 := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP1(t)
		metrolinkDeparturesForAtcoCode9400ZZMASTP2 := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP2(t)
		metrolinkDeparturesForAtcoCode9400ZZMASTP3 := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP3(t)
		metrolinkDeparturesForAtcoCode9400ZZMASTP4 := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP4(t)
		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP").Return(nil, redis.ErrNil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP1").Return(metrolinkDeparturesForAtcoCode9400ZZMASTP1, nil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP2").Return(metrolinkDeparturesForAtcoCode9400ZZMASTP2, nil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP3").Return(metrolinkDeparturesForAtcoCode9400ZZMASTP3, nil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP4").Return(metrolinkDeparturesForAtcoCode9400ZZMASTP4, nil)

		lastUpdatedTime := givenLastUpdatedTime(t)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(lastUpdatedTime, nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkStopAreaCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, thenExpectJsonDeparturesWithPlatform(t, validMetrolinkStopAreaCode, thenExpectDeparturesFor940GZZMASTP(t), lastUpdatedTime), readJson(t, rc))
	})

	t.Run(`Given a valid Metrolink AtcoCode is requested
And there are no departures for that AtcoCode
When Json is called
Then the departures value is an empty slice`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)

		validMetrolinkAtcoCode := "9400ZZMASTP1"

		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, validMetrolinkAtcoCode).Return(nil, nil)

		lastUpdatedTime := givenLastUpdatedTime(t)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(lastUpdatedTime, nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkAtcoCode)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, thenExpectJsonWithEmptyDeparturesSlice(t, validMetrolinkAtcoCode, lastUpdatedTime), readJson(t, rc))
	})

	t.Run(`Given an error occurs getting AtcoCodes for a StopAreaCode
When Json is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		validMetrolinkStopAreaCode := "940GZZMASTP"

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)
		stopsInAreaGetterErr := errors.New("FUBAR")
		stopsInAreaGetter.EXPECT().GetStopsInArea(ctx, validMetrolinkStopAreaCode).Return(nil, stopsInAreaGetterErr)

		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(givenLastUpdatedTime(t), nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkStopAreaCode)

		// Then
		assert.Nil(t, rc)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error getting ATCO codes for '940GZZMASTP': FUBAR")
	})

	t.Run(`Given an error occurs fetching departures for an AtcoCode
And the error is not redis.ErrNil
When Json is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		validMetrolinkStopAreaCode := "940GZZMASTP"

		stopsInAreaGetter := mock_repository.NewMockStopsInAreaGetter(ctrl)
		stopsInAreaGetter.EXPECT().GetStopsInArea(ctx, validMetrolinkStopAreaCode).Return([]string{"9400ZZMASTP", "9400ZZMASTP1", "9400ZZMASTP2", "9400ZZMASTP3", "9400ZZMASTP4"}, nil)

		metrolinkDeparturesForAtcoCode9400ZZMASTP1 := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP1(t)
		metrolinkDeparturesForAtcoCode9400ZZMASTP2 := givenMetrolinkDeparturesForAtcoCode9400ZZMASTP2(t)
		metrolinkDeparturesGetter := mock_repository.NewMockMetrolinkDeparturesGetter(ctrl)
		metrolinkDeparturesGetterErr := errors.New("FUBAR")
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP").Return(nil, redis.ErrNil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP1").Return(metrolinkDeparturesForAtcoCode9400ZZMASTP1, nil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP2").Return(metrolinkDeparturesForAtcoCode9400ZZMASTP2, nil)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP3").Return(nil, metrolinkDeparturesGetterErr)
		metrolinkDeparturesGetter.EXPECT().Get(ctx, "9400ZZMASTP4").Return(nil, metrolinkDeparturesGetterErr)

		metrolinkDeparturesSystemStatusGetter := mock_repository.NewMockSystemStatusGetter(ctrl)
		metrolinkDeparturesSystemStatusGetter.EXPECT().Get(ctx).Return(givenLastUpdatedTime(t), nil)

		api := NewApi(logger, stopsInAreaGetter, metrolinkDeparturesGetter, metrolinkDeparturesSystemStatusGetter, givenCurrentTimeFunc(t), givenStaleDataThreshold(t), givenTimeLocation(t))

		// When
		rc, statusCode, err := api.Json(ctx, validMetrolinkStopAreaCode)

		// Then
		assert.Nil(t, rc)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error fetching Metrolink departures for '940GZZMASTP': 2 errors occurred:\n\t")
		assert.Contains(t, err.Error(), "* error getting departures for AtcoCode '9400ZZMASTP3': FUBAR\n")
		assert.Contains(t, err.Error(), "* error getting departures for AtcoCode '9400ZZMASTP4': FUBAR\n")
	})
}
