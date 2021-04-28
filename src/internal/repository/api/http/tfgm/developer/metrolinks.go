package developer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type MetrolinkDepartures struct {
	PassengerInformationDisplays []*PassengerInformationDisplay `json:"value"`
}

func (md *MetrolinkDepartures) LastUpdated() time.Time {
	lastUpdated := time.Time{}

	for _, passengerInformationDisplay := range md.PassengerInformationDisplays {
		if passengerInformationDisplay.LastUpdated.After(lastUpdated) {
			lastUpdated = passengerInformationDisplay.LastUpdated
		}
	}

	return lastUpdated
}

type PassengerInformationDisplay struct {
	Id              int
	Line            string
	TLAREF          string
	PIDREF          string
	StationLocation string
	AtcoCode        string
	Direction       string
	Dest0           string
	Carriages0      string
	Status0         string
	Wait0           string
	Dest1           string
	Carriages1      string
	Status1         string
	Wait1           string
	Dest2           string
	Carriages2      string
	Status2         string
	Wait2           string
	Dest3           string
	Carriages3      string
	Status3         string
	Wait3           string
	MessageBoard    string
	LastUpdated     time.Time
}

type TfgmDeveloperMetrolinkDataSource struct {
	logger     *zap.Logger
	httpClient *http.Client
	url        string
	apiKey     string
}

func NewTfgmDeveloperMetrolinkDataSource(logger *zap.Logger, httpClient *http.Client, url string, apiKey string) *TfgmDeveloperMetrolinkDataSource {
	return &TfgmDeveloperMetrolinkDataSource{
		logger:     logger,
		httpClient: httpClient,
		url:        url,
		apiKey:     apiKey,
	}
}

func (ds *TfgmDeveloperMetrolinkDataSource) Fetch(ctx context.Context) (*domain.MetrolinkDepartures, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", ds.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", ds.apiKey)

	resp, err := ds.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)

		errMsg := "error response from data source"

		ds.logger.Error(errMsg, zap.Int("StatusCode", resp.StatusCode), zap.String("Status", resp.Status), zap.String("Body", buf.String()))

		return nil, fmt.Errorf("%s: %s", errMsg, resp.Status)
	}

	var metrolinkDepartures MetrolinkDepartures
	if err := json.NewDecoder(resp.Body).Decode(&metrolinkDepartures); err != nil {
		return nil, errors.Wrap(err, "error decoding body as JSON")
	}

	if len(metrolinkDepartures.PassengerInformationDisplays) == 0 {
		return nil, errors.New("no departures data returned from data source")
	}

	for i := range metrolinkDepartures.PassengerInformationDisplays {
		metrolinkDepartures.PassengerInformationDisplays[i].LastUpdated = ds.correctLastUpdatedTime(metrolinkDepartures.PassengerInformationDisplays[i].LastUpdated)
	}

	return &domain.MetrolinkDepartures{
		Departures:  ds.convertToDomainMetrolinkDepartures(&metrolinkDepartures),
		LastUpdated: metrolinkDepartures.LastUpdated(),
	}, nil
}

func (ds *TfgmDeveloperMetrolinkDataSource) convertToDomainMetrolinkDepartures(metrolinkDepartures *MetrolinkDepartures) []*domain.MetrolinkDeparture {
	var domainMetrolinkDepartures []*domain.MetrolinkDeparture

	for _, passengerInformationDisplay := range ds.filterDuplicatePassengerInformationDisplays(metrolinkDepartures.PassengerInformationDisplays) {
		if passengerInformationDisplay.Status0 != "" {
			domainMetrolinkDepartures = append(domainMetrolinkDepartures, &domain.MetrolinkDeparture{
				AtcoCode:    passengerInformationDisplay.AtcoCode,
				Order:       0,
				Destination: passengerInformationDisplay.Dest0,
				Carriages:   passengerInformationDisplay.Carriages0,
				Status:      passengerInformationDisplay.Status0,
				Wait:        passengerInformationDisplay.Wait0,
				LastUpdated: passengerInformationDisplay.LastUpdated,
			})
		}

		if passengerInformationDisplay.Status1 != "" {
			domainMetrolinkDepartures = append(domainMetrolinkDepartures, &domain.MetrolinkDeparture{
				AtcoCode:    passengerInformationDisplay.AtcoCode,
				Order:       1,
				Destination: passengerInformationDisplay.Dest1,
				Carriages:   passengerInformationDisplay.Carriages1,
				Status:      passengerInformationDisplay.Status1,
				Wait:        passengerInformationDisplay.Wait1,
				LastUpdated: passengerInformationDisplay.LastUpdated,
			})
		}

		if passengerInformationDisplay.Status2 != "" {
			domainMetrolinkDepartures = append(domainMetrolinkDepartures, &domain.MetrolinkDeparture{
				AtcoCode:    passengerInformationDisplay.AtcoCode,
				Order:       2,
				Destination: passengerInformationDisplay.Dest2,
				Carriages:   passengerInformationDisplay.Carriages2,
				Status:      passengerInformationDisplay.Status2,
				Wait:        passengerInformationDisplay.Wait2,
				LastUpdated: passengerInformationDisplay.LastUpdated,
			})
		}

		if passengerInformationDisplay.Status3 != "" {
			domainMetrolinkDepartures = append(domainMetrolinkDepartures, &domain.MetrolinkDeparture{
				AtcoCode:    passengerInformationDisplay.AtcoCode,
				Order:       3,
				Destination: passengerInformationDisplay.Dest3,
				Carriages:   passengerInformationDisplay.Carriages3,
				Status:      passengerInformationDisplay.Status3,
				Wait:        passengerInformationDisplay.Wait3,
				LastUpdated: passengerInformationDisplay.LastUpdated,
			})
		}
	}

	return domainMetrolinkDepartures
}

func (ds *TfgmDeveloperMetrolinkDataSource) filterDuplicatePassengerInformationDisplays(passengerInformationDisplays []*PassengerInformationDisplay) []*PassengerInformationDisplay {
	var x struct{}

	processedAtcoCodes := make(map[string]*struct{})

	n := 0
	for _, passengerInformationDisplay := range passengerInformationDisplays {
		if processedAtcoCodes[passengerInformationDisplay.AtcoCode] == nil {
			passengerInformationDisplays[n] = passengerInformationDisplay
			processedAtcoCodes[passengerInformationDisplay.AtcoCode] = &x
			n++
		}
	}

	return passengerInformationDisplays[:n]
}

// The Metrolinks API currently returns the LastUpdated time values as local time, but with a "Z" suffix
// This means that the LastUpdated time value is technically incorrect:
// e.g. a request made at 2021-04-24T15:04:05+01:00 (British Summer Time)
// returns a LastUpdated value of 2021-04-24T15:04:05Z
// this is equal to 2021-04-24T16:04:05+01:00
// The API data suggests it was last updated one hour in the future!
// This function corrects this error.
func (ds *TfgmDeveloperMetrolinkDataSource) correctLastUpdatedTime(lastUpdated time.Time) time.Time {
	europeLondon, _ := time.LoadLocation("Europe/London")

	lastUpdatedInEuropeLondon := time.Date(lastUpdated.Year(), lastUpdated.Month(), lastUpdated.Day(), lastUpdated.Hour(), lastUpdated.Minute(), lastUpdated.Second(), lastUpdated.Nanosecond(), europeLondon)

	if lastUpdatedInEuropeLondon.Equal(lastUpdated.UTC()) {
		return lastUpdated
	}

	return lastUpdated.Add(-time.Hour)
}
