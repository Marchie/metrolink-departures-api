package developer_test

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/api/http/tfgm/developer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mockLogger(t *testing.T) *zap.Logger {
	t.Helper()

	zapCore, _ := observer.New(zapcore.DebugLevel)
	return zap.New(zapCore)
}

func mockMetrolinksWorkingServer(t *testing.T, expApiKey string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Ocp-Apim-Subscription-Key") != expApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{ "status": "not authorized" }`))
			return
		}

		body := `{
    "@odata.context": "https://opendataclientapi.azurewebsites.net/odata/$metadata#Metrolinks",
    "value": [
        {
            "Id": 846,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID05",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP1",
            "Direction": "Incoming",
            "Dest0": "Victoria",
            "Carriages0": "Single",
            "Status0": "Due",
            "Wait0": "2",
            "Dest1": "Ashton-under-Lyne",
            "Carriages1": "Double",
            "Status1": "Due",
            "Wait1": "9",
            "Dest2": "",
            "Carriages2": "",
            "Status2": "",
            "Wait2": "",
            "Dest3": "Piccadilly",
            "Carriages3": "Double",
            "Status3": "Due",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "12",
            "LastUpdated": "2021-03-21T15:34:54Z"
        },
        {
            "Id": 842,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID01",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP4",
            "Direction": "Outgoing",
            "Dest0": "Altrincham",
            "Carriages0": "Double",
            "Status0": "Due",
            "Wait0": "6",
            "Dest1": "Manchester Airport",
            "Carriages1": "Single",
            "Status1": "Due",
            "Wait1": "10",
            "Dest2": "MediaCityUK",
            "Carriages2": "Double",
            "Status2": "Due",
            "Wait2": "14",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "Please see printed posters forfirst and last tram times.info: www.metrolink.co.uk",
            "Wait3": "",
            "LastUpdated": "2021-03-21T15:34:53Z"
        },
        {
            "Id": 844,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID03",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP4",
            "Direction": "Outgoing",
            "Dest0": "Altrincham",
            "Carriages0": "Double",
            "Status0": "Due",
            "Wait0": "6",
            "Dest1": "Manchester Airport",
            "Carriages1": "Single",
            "Status1": "Due",
            "Wait1": "10",
            "Dest2": "MediaCityUK",
            "Carriages2": "Double",
            "Status2": "Due",
            "Wait2": "14",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "",
            "LastUpdated": "2021-03-21T15:34:53Z"
        }
    ]
}`

		w.Header().Add("Content-type", "application/json; odata.metadata=minimal")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func mockMetrolinksWorkingServerInBritishSummerTime(t *testing.T, expApiKey string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Ocp-Apim-Subscription-Key") != expApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{ "status": "not authorized" }`))
			return
		}

		body := `{
    "@odata.context": "https://opendataclientapi.azurewebsites.net/odata/$metadata#Metrolinks",
    "value": [
        {
            "Id": 846,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID05",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP1",
            "Direction": "Incoming",
            "Dest0": "Victoria",
            "Carriages0": "Single",
            "Status0": "Due",
            "Wait0": "2",
            "Dest1": "Ashton-under-Lyne",
            "Carriages1": "Double",
            "Status1": "Due",
            "Wait1": "9",
            "Dest2": "",
            "Carriages2": "",
            "Status2": "",
            "Wait2": "",
            "Dest3": "Piccadilly",
            "Carriages3": "Double",
            "Status3": "Due",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "12",
            "LastUpdated": "2021-04-21T15:34:54Z"
        },
        {
            "Id": 842,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID01",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP4",
            "Direction": "Outgoing",
            "Dest0": "Altrincham",
            "Carriages0": "Double",
            "Status0": "Due",
            "Wait0": "6",
            "Dest1": "Manchester Airport",
            "Carriages1": "Single",
            "Status1": "Due",
            "Wait1": "10",
            "Dest2": "MediaCityUK",
            "Carriages2": "Double",
            "Status2": "Due",
            "Wait2": "14",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "Please see printed posters forfirst and last tram times.info: www.metrolink.co.uk",
            "Wait3": "",
            "LastUpdated": "2021-04-21T15:34:53Z"
        },
        {
            "Id": 844,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID03",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP4",
            "Direction": "Outgoing",
            "Dest0": "Altrincham",
            "Carriages0": "Double",
            "Status0": "Due",
            "Wait0": "6",
            "Dest1": "Manchester Airport",
            "Carriages1": "Single",
            "Status1": "Due",
            "Wait1": "10",
            "Dest2": "MediaCityUK",
            "Carriages2": "Double",
            "Status2": "Due",
            "Wait2": "14",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "",
            "LastUpdated": "2021-04-21T15:34:53Z"
        }
    ]
}`

		w.Header().Add("Content-type", "application/json; odata.metadata=minimal")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func mockMetrolinksWorkingServerNoDepartures(t *testing.T, expApiKey string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Ocp-Apim-Subscription-Key") != expApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{ "status": "not authorized" }`))
			return
		}

		body := `{
    "@odata.context": "https://opendataclientapi.azurewebsites.net/odata/$metadata#Metrolinks",
    "value": [
        {
            "Id": 846,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID05",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP1",
            "Direction": "Incoming",
            "Dest0": "",
            "Carriages0": "",
            "Status0": "",
            "Wait0": "",
            "Dest1": "",
            "Carriages1": "",
            "Status1": "",
            "Wait1": "",
            "Dest2": "",
            "Carriages2": "",
            "Status2": "",
            "Wait2": "",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "",
            "LastUpdated": "2021-03-21T15:34:54Z"
        },
        {
            "Id": 842,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID01",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP4",
            "Direction": "Outgoing",
			"Dest0": "",
            "Carriages0": "",
            "Status0": "",
            "Wait0": "",
            "Dest1": "",
            "Carriages1": "",
            "Status1": "",
            "Wait1": "",
            "Dest2": "",
            "Carriages2": "",
            "Status2": "",
            "Wait2": "",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "",
            "LastUpdated": "2021-03-21T15:34:54Z"
        },
        {
            "Id": 844,
            "Line": "Eccles",
            "TLAREF": "SPS",
            "PIDREF": "SPS-PID03",
            "StationLocation": "St Peter's Square",
            "AtcoCode": "9400ZZMASTP4",
            "Direction": "Outgoing",
            "Dest0": "",
            "Carriages0": "",
            "Status0": "",
            "Wait0": "",
            "Dest1": "",
            "Carriages1": "",
            "Status1": "",
            "Wait1": "",
            "Dest2": "",
            "Carriages2": "",
            "Status2": "",
            "Wait2": "",
            "Dest3": "",
            "Carriages3": "",
            "Status3": "",
            "MessageBoard": "PLANNED IMPROVEMENT WORKS - No service between Eccles and MediaCityUK. A bus replacement service is operating at Eccles and MediaCityUK only. For more info please visit tfgm.com",
            "Wait3": "",
            "LastUpdated": "2021-03-21T15:34:54Z"
        }
    ]
}`

		w.Header().Add("Content-type", "application/json; odata.metadata=minimal")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func mockMetrolinksNoData(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := `{
	"@odata.context": "https://opendataclientapi.azurewebsites.net/odata/$metadata#Metrolinks",
    "value": []
}`

		w.Header().Add("Content-type", "application/json; odata.metadata=minimal")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func mockMetrolinksInvalidJson(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := `X`

		w.Header().Add("Content-type", "application/json; odata.metadata=minimal")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func mockMetrolinksTeapot(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
}

func TestMetrolinkDepartures_LastUpdated(t *testing.T) {
	t.Run(`Given MetrolinkDepartures data
When LastUpdated() is called
Then the latest lastUpdated value from the data is returned`, func(t *testing.T) {
		// Given
		metrolinkDepartures := developer.MetrolinkDepartures{
			PassengerInformationDisplays: []*developer.PassengerInformationDisplay{
				{
					LastUpdated: time.Date(2021, time.March, 27, 2, 21, 0, 0, time.UTC),
				},
				{
					LastUpdated: time.Date(2021, time.March, 27, 2, 21, 27, 0, time.UTC),
				},
				{
					LastUpdated: time.Date(2021, time.March, 27, 2, 21, 13, 0, time.UTC),
				},
			},
		}

		// When
		result := metrolinkDepartures.LastUpdated()

		// Then
		expLastUpdated := time.Date(2021, time.March, 27, 2, 21, 27, 0, time.UTC)

		assert.Equal(t, expLastUpdated, result)
	})
}

func TestTfgmDeveloperMetrolinkDataSource_Fetch(t *testing.T) {
	t.Run(`Given data is available from the TfGM Developer Metrolinks API
When data is fetched from the API
Then a slice containing all departures is returned
And duplicated information caused by there being multiple Passenger Information Displays for an AtcoCode is removed`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksWorkingServer(t, apiKey)
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		result, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.Nil(t, err)

		expLastUpdated, err := time.Parse(time.RFC3339, "2021-03-21T15:34:54Z")
		if err != nil {
			t.Fatal(err)
		}

		expDomainMetrolinkDepartures := &domain.MetrolinkDepartures{
			Departures: []*domain.MetrolinkDeparture{
				{
					AtcoCode:    "9400ZZMASTP1",
					Order:       0,
					Destination: "Victoria",
					Carriages:   "Single",
					Status:      "Due",
					Wait:        "2",
					LastUpdated: expLastUpdated,
				},
				{
					AtcoCode:    "9400ZZMASTP1",
					Order:       1,
					Destination: "Ashton-under-Lyne",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "9",
					LastUpdated: expLastUpdated,
				},
				{
					AtcoCode:    "9400ZZMASTP1",
					Order:       3,
					Destination: "Piccadilly",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "12",
					LastUpdated: expLastUpdated,
				},
				{
					AtcoCode:    "9400ZZMASTP4",
					Order:       0,
					Destination: "Altrincham",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "6",
					LastUpdated: expLastUpdated.Add(-time.Second),
				},
				{
					AtcoCode:    "9400ZZMASTP4",
					Order:       1,
					Destination: "Manchester Airport",
					Carriages:   "Single",
					Status:      "Due",
					Wait:        "10",
					LastUpdated: expLastUpdated.Add(-time.Second),
				},
				{
					AtcoCode:    "9400ZZMASTP4",
					Order:       2,
					Destination: "MediaCityUK",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "14",
					LastUpdated: expLastUpdated.Add(-time.Second),
				},
			},
			LastUpdated: expLastUpdated,
		}

		assert.EqualValues(t, expDomainMetrolinkDepartures, result)
	})

	t.Run(`Given the TfGM Developer Metrolinks API is available
And no trams are running
When data is fetched from the API
Then an empty departures slice is returned
And the LastUpdated time is populated`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksWorkingServerNoDepartures(t, apiKey)
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		result, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.Nil(t, err)

		expLastUpdated, err := time.Parse(time.RFC3339, "2021-03-21T15:34:54Z")
		if err != nil {
			t.Fatal(err)
		}

		expDomainMetrolinkDepartures := &domain.MetrolinkDepartures{
			Departures:  nil,
			LastUpdated: expLastUpdated,
		}

		assert.EqualValues(t, expDomainMetrolinkDepartures, result)
	})

	t.Run(`Given data is available from the TfGM Developer Metrolinks API
And Daylight Savings Time is in effect
When data is fetched from the API
And the API returns LastUpdated values an hour into the future
Then a slice containing all departures is returned
And the LastUpdated times are corrected to their true values`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksWorkingServerInBritishSummerTime(t, apiKey)
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		result, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.Nil(t, err)

		expLastUpdated, err := time.Parse(time.RFC3339, "2021-04-21T14:34:54Z")
		if err != nil {
			t.Fatal(err)
		}

		expDomainMetrolinkDepartures := &domain.MetrolinkDepartures{
			Departures: []*domain.MetrolinkDeparture{
				{
					AtcoCode:    "9400ZZMASTP1",
					Order:       0,
					Destination: "Victoria",
					Carriages:   "Single",
					Status:      "Due",
					Wait:        "2",
					LastUpdated: expLastUpdated,
				},
				{
					AtcoCode:    "9400ZZMASTP1",
					Order:       1,
					Destination: "Ashton-under-Lyne",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "9",
					LastUpdated: expLastUpdated,
				},
				{
					AtcoCode:    "9400ZZMASTP1",
					Order:       3,
					Destination: "Piccadilly",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "12",
					LastUpdated: expLastUpdated,
				},
				{
					AtcoCode:    "9400ZZMASTP4",
					Order:       0,
					Destination: "Altrincham",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "6",
					LastUpdated: expLastUpdated.Add(-time.Second),
				},
				{
					AtcoCode:    "9400ZZMASTP4",
					Order:       1,
					Destination: "Manchester Airport",
					Carriages:   "Single",
					Status:      "Due",
					Wait:        "10",
					LastUpdated: expLastUpdated.Add(-time.Second),
				},
				{
					AtcoCode:    "9400ZZMASTP4",
					Order:       2,
					Destination: "MediaCityUK",
					Carriages:   "Double",
					Status:      "Due",
					Wait:        "14",
					LastUpdated: expLastUpdated.Add(-time.Second),
				},
			},
			LastUpdated: expLastUpdated,
		}

		assert.EqualValues(t, expDomainMetrolinkDepartures, result)
	})

	t.Run(`Given the URL provided for the TfGM Developer Metrolinks API cannot be reached
When data is fetched from the TfGM Developer Metrolinks API
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServerUrl := "http://fail"

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServerUrl, apiKey)

		ctx := context.Background()

		// When
		_, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.NotNil(t, err)

		assert.EqualError(t, err, `Get "http://fail": dial tcp: lookup fail: no such host`)
	})

	t.Run(`Given the URL provided for the TfGM Developer Metrolinks API returns a non-200 status
When data is fetched from the TfGM Developer Metrolinks API
Then a descriptive error is logged
And an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		zapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		logger := zap.New(zapCore)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksWorkingServer(t, "invalid")
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		_, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.NotNil(t, err)

		assert.EqualError(t, err, "error response from data source: 401 Unauthorized")

		assert.Equal(t, observedLogs.Len(), 1)

		logs := observedLogs.TakeAll()

		assert.Equal(t, "error response from data source", logs[0].Message)
		assert.Equal(t, int64(http.StatusUnauthorized), logs[0].ContextMap()["StatusCode"])
		assert.Equal(t, "401 Unauthorized", logs[0].ContextMap()["Status"])
		assert.Equal(t, `{ "status": "not authorized" }`, logs[0].ContextMap()["Body"])
	})

	t.Run(`Given the URL provided for the TfGM Developer Metrolinks API returns a non-200 status
And the TfGM Developer Metrolinks API response does not include a body
When data is fetched from the TfGM Developer Metrolinks API
Then a descriptive error is logged
And an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		zapCore, observedLogs := observer.New(zapcore.ErrorLevel)
		logger := zap.New(zapCore)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksTeapot(t)
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		_, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.NotNil(t, err)

		assert.EqualError(t, err, "error response from data source: 418 I'm a teapot")

		assert.Equal(t, observedLogs.Len(), 1)

		logs := observedLogs.TakeAll()

		assert.Equal(t, "error response from data source", logs[0].Message)
		assert.Equal(t, int64(http.StatusTeapot), logs[0].ContextMap()["StatusCode"])
		assert.Equal(t, "418 I'm a teapot", logs[0].ContextMap()["Status"])
		assert.Equal(t, "", logs[0].ContextMap()["Body"])
	})

	t.Run(`Given the TfGM Developer Metrolinks API returns invalid JSON
When data is fetched from the TfGM Developer Metrolinks API
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksInvalidJson(t)
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		_, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.NotNil(t, err)

		assert.EqualError(t, err, "error decoding body as JSON: invalid character 'X' looking for beginning of value")
	})

	t.Run(`Given the TfGM Developer Metrolinks API returns an empty value array
When data is fetched from the TfGM Developer Metrolinks API
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := mockLogger(t)

		httpClient := &http.Client{}

		apiKey := "abc123"

		metrolinksServer := mockMetrolinksNoData(t)
		defer metrolinksServer.Close()

		tfgmDeveloperMetrolinkDataSource := developer.NewTfgmDeveloperMetrolinkDataSource(logger, httpClient, metrolinksServer.URL, apiKey)

		ctx := context.Background()

		// When
		result, err := tfgmDeveloperMetrolinkDataSource.Fetch(ctx)

		// Then
		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}
