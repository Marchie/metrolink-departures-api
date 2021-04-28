package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/domain"
	"github.com/Marchie/tf-experiment/lambda/internal/repository"
	"github.com/Marchie/tf-experiment/lambda/pkg/tfgm"
	"github.com/gomodule/redigo/redis"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Api struct {
	logger                    *zap.Logger
	stopsInAreaGetter         repository.StopsInAreaGetter
	metrolinkDeparturesGetter repository.MetrolinkDeparturesGetter
	systemStatusGetter        repository.SystemStatusGetter
	currentTimeFunc           func() time.Time
	staleDataThreshold        time.Duration
	timeLocation              *time.Location
}

func NewApi(logger *zap.Logger, stopsInAreaGetter repository.StopsInAreaGetter, metrolinkDeparturesGetter repository.MetrolinkDeparturesGetter, systemStatusGetter repository.SystemStatusGetter, currentTimeFunc func() time.Time, staleDataThreshold time.Duration, timeLocation *time.Location) *Api {
	return &Api{
		logger:                    logger,
		stopsInAreaGetter:         stopsInAreaGetter,
		metrolinkDeparturesGetter: metrolinkDeparturesGetter,
		systemStatusGetter:        systemStatusGetter,
		currentTimeFunc:           currentTimeFunc,
		staleDataThreshold:        staleDataThreshold,
		timeLocation:              timeLocation,
	}
}

func (m *Api) Json(ctx context.Context, stopAreaCodeOrAtcoCode string) (io.ReadCloser, int, error) {
	stopAreaCodeOrAtcoCode = strings.ToUpper(stopAreaCodeOrAtcoCode)

	if !m.validateStopAreaCodeOrAtcoCode(stopAreaCodeOrAtcoCode) {
		return m.encodeJsonErrorResponse(stopAreaCodeOrAtcoCode, http.StatusBadRequest, "invalid StopAreaCode or AtcoCode")
	}

	lastUpdated, err := m.systemStatusGetter.Get(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrapf(err, "error getting Metrolink departures system status")
	}

	if m.currentTimeFunc().Sub(*lastUpdated) > m.staleDataThreshold {
		return m.encodeJsonErrorResponse(stopAreaCodeOrAtcoCode, http.StatusBadGateway, fmt.Sprintf("Metrolink departures data is outdated: last updated at %s", lastUpdated.Format(time.RFC3339)))
	}

	atcoCodes, err := m.atcoCodesToQuery(ctx, stopAreaCodeOrAtcoCode)
	if err != nil {
		if strings.HasSuffix(err.Error(), redis.ErrNil.Error()) {
			return m.encodeJsonErrorResponse(stopAreaCodeOrAtcoCode, http.StatusBadRequest, "invalid StopAreaCode")
		}

		return nil, http.StatusInternalServerError, errors.Wrapf(err, "error getting ATCO codes for '%s'", stopAreaCodeOrAtcoCode)
	}

	departures, err := m.getDeparturesForAtcoCodes(ctx, atcoCodes)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrapf(err, "error fetching Metrolink departures for '%s'", stopAreaCodeOrAtcoCode)
	}

	sort.Sort(ByWaitStatusDestinationOrderPlatformCarriages(departures))

	return m.encodeJsonSuccessResponse(stopAreaCodeOrAtcoCode, departures, *lastUpdated)
}

func (m *Api) validateStopAreaCodeOrAtcoCode(stopAreaCodeOrAtcoCode string) bool {
	match, _ := regexp.MatchString("^940[0G]ZZMA[A-Z]{3}[1-4]?$", stopAreaCodeOrAtcoCode)
	return match
}

func (m *Api) atcoCodesToQuery(ctx context.Context, stopAreaCodeOrAtcoCode string) ([]string, error) {
	if strings.HasPrefix(stopAreaCodeOrAtcoCode, "9400") {
		return []string{stopAreaCodeOrAtcoCode}, nil
	}

	return m.stopsInAreaGetter.GetStopsInArea(ctx, stopAreaCodeOrAtcoCode)
}

func (m *Api) getDeparturesForAtcoCodes(ctx context.Context, atcoCodes []string) ([]*domain.MetrolinkDeparture, error) {
	wg := sync.WaitGroup{}
	chMetrolinkDepartures := make(chan []*domain.MetrolinkDeparture)
	chErr := make(chan error)

	for _, atcoCode := range atcoCodes {
		wg.Add(1)
		go m.getDeparturesForAtcoCode(&wg, ctx, atcoCode, chMetrolinkDepartures, chErr)
	}

	chCombinedMetrolinkDepartures := m.combineMetrolinkDepartures(chMetrolinkDepartures)
	chCombinedErr := m.combineErrors(chErr)

	go func() {
		defer close(chErr)
		defer close(chMetrolinkDepartures)

		wg.Wait()
	}()

	return <-chCombinedMetrolinkDepartures, <-chCombinedErr
}

func (m *Api) getDeparturesForAtcoCode(wg *sync.WaitGroup, ctx context.Context, atcoCode string, chMetrolinkDepartures chan []*domain.MetrolinkDeparture, chErr chan error) {
	defer wg.Done()

	metrolinkDepartures, err := m.metrolinkDeparturesGetter.Get(ctx, atcoCode)
	if err != nil {
		if err == redis.ErrNil {
			return
		}

		chErr <- errors.Wrapf(err, "error getting departures for AtcoCode '%s'", atcoCode)
		return
	}

	chMetrolinkDepartures <- metrolinkDepartures
}

func (m *Api) combineMetrolinkDepartures(chMetrolinkDepartures <-chan []*domain.MetrolinkDeparture) chan []*domain.MetrolinkDeparture {
	chCombinedMetrolinkDepartures := make(chan []*domain.MetrolinkDeparture, 1)

	go func() {
		defer close(chCombinedMetrolinkDepartures)

		var combinedMetrolinkDepartures []*domain.MetrolinkDeparture

		for metrolinkDepartures := range chMetrolinkDepartures {
			combinedMetrolinkDepartures = append(combinedMetrolinkDepartures, metrolinkDepartures...)
		}

		chCombinedMetrolinkDepartures <- combinedMetrolinkDepartures
	}()

	return chCombinedMetrolinkDepartures
}

func (m *Api) combineErrors(chErr <-chan error) chan error {
	chCombinedErr := make(chan error, 1)

	go func() {
		defer close(chCombinedErr)

		var combinedErr error

		for err := range chErr {
			combinedErr = multierror.Append(combinedErr, err)
		}

		chCombinedErr <- combinedErr
	}()

	return chCombinedErr
}

func (m *Api) convertToPublicApi(stopAreaCodeOrAtcoCode string, departures []*domain.MetrolinkDeparture, lastUpdated time.Time) *tfgm.MetrolinkDepartures {
	convertedDepartures := make([]*tfgm.MetrolinkDeparture, 0)

	for sequence, departure := range departures {
		convertedDepartures = append(convertedDepartures, &tfgm.MetrolinkDeparture{
			AtcoCode:    departure.AtcoCode,
			Sequence:    sequence,
			Destination: departure.Destination,
			Status:      departure.Status,
			Wait:        departure.Wait,
			Carriages:   departure.Carriages,
			Platform:    departure.Platform,
			LastUpdated: departure.LastUpdated.In(m.timeLocation),
		})
	}

	return &tfgm.MetrolinkDepartures{
		RequestedLocation: stopAreaCodeOrAtcoCode,
		Departures:        convertedDepartures,
		LastUpdated:       lastUpdated.In(m.timeLocation),
	}
}

func (m *Api) encodeJsonResponse(v interface{}, statusCode int) (io.ReadCloser, int, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")

	if err := enc.Encode(v); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return ioutil.NopCloser(buf), statusCode, nil
}

func (m *Api) encodeJsonSuccessResponse(stopAreaCodeOrAtcoCode string, departures []*domain.MetrolinkDeparture, lastUpdated time.Time) (io.ReadCloser, int, error) {
	return m.encodeJsonResponse(m.convertToPublicApi(stopAreaCodeOrAtcoCode, departures, lastUpdated), http.StatusOK)
}

func (m *Api) encodeJsonErrorResponse(stopAreaCodeOrAtcoCode string, statusCode int, errorMsg string) (io.ReadCloser, int, error) {
	return m.encodeJsonResponse(map[string]string{
		"requestedLocation": stopAreaCodeOrAtcoCode,
		"error":             errorMsg,
	}, statusCode)
}
