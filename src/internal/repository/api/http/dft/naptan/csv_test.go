package naptan_test

import (
	"bytes"
	"context"
	mock_http "github.com/Marchie/tf-experiment/lambda/internal/mocks/repository/api/http"
	mock_compression "github.com/Marchie/tf-experiment/lambda/internal/mocks/repository/compression"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/api/http/dft/naptan"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func givenArbitraryZipData(t *testing.T, data []byte) io.ReadCloser {
	t.Helper()

	return ioutil.NopCloser(bytes.NewBuffer(data))
}

func givenStopsInAreaCsv(t *testing.T) io.ReadCloser {
	t.Helper()

	csvData := strings.Join([]string{
		"StopAreaCode,AtcoCode,CreationDateTime,ModificationDateTime,RevisionNumber,Modification",
		"940GZZMAAUD,9400ZZMAAUD2,2013-09-03T09:52:00,2013-09-03T09:52:00,0,new",
		"940GZZMAAWT,9400ZZMAAWT1,2013-09-04T10:29:00,2013-09-04T10:29:00,0,new",
		"940GZZMASTP,9400ZZMASTP4,,2020-02-26T11:56:45,1,rev",
		"940GZZMASTP,9400ZZMASTP2,2006-12-12T00:00:00,2006-12-12T00:00:00,0,new",
		"940GZZMASTP,9400ZZMASTP,2010-03-01T14:08:00,2010-03-01T14:08:00,0,new",
		"940GZZMASTP,9400ZZMASTP3,,2020-02-26T11:56:45,1,rev",
		"940GZZMASTP,9400ZZMASTP1,2006-12-12T00:00:00,2006-12-12T00:00:00,0,new",
	}, "\n")

	return ioutil.NopCloser(bytes.NewBufferString(csvData))
}

func givenEmptyCsvData(t *testing.T) io.ReadCloser {
	t.Helper()

	return ioutil.NopCloser(bytes.NewBufferString(""))
}

func givenCorruptCsvData(t *testing.T) io.ReadCloser {
	t.Helper()

	return ioutil.NopCloser(bytes.NewBufferString("not CSV data\nanother,line\n\n\n"))
}

func TestCSV_FetchStopsInArea(t *testing.T) {
	t.Run(`Given valid NaPTAN CSV data is retrieved
When FetchStopsInArea is called
Then a map of StopAreaCodes to StopAreas is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		zipData := []byte("arbitrary data")

		stopsInAreaFilename := "StopsInArea.csv"
		stopAreaCodeColumnIndex := 0
		atcoCodeColumnIndex := 1

		zipFileFetcher := mock_http.NewMockZipFileFetcher(ctrl)
		zipFileFetcher.EXPECT().FetchZipFile(ctx).Return(givenArbitraryZipData(t, zipData), nil)

		extractor := mock_compression.NewMockExtractor(ctrl)
		extractor.EXPECT().ExtractFile(zipData, stopsInAreaFilename).Return(givenStopsInAreaCsv(t), nil)

		naptanCsv := naptan.NewCSV(logger, zipFileFetcher, extractor, stopsInAreaFilename, stopAreaCodeColumnIndex, atcoCodeColumnIndex)

		// When
		stopsInAreaMap, err := naptanCsv.FetchStopsInArea(ctx)

		// Then
		assert.NotNil(t, stopsInAreaMap)
		assert.Nil(t, err)

		expectedStopsInAreaMap := make(map[string][]string)
		expectedStopsInAreaMap["940GZZMAAUD"] = []string{"9400ZZMAAUD2"}
		expectedStopsInAreaMap["940GZZMAAWT"] = []string{"9400ZZMAAWT1"}
		expectedStopsInAreaMap["940GZZMASTP"] = []string{"9400ZZMASTP4", "9400ZZMASTP2", "9400ZZMASTP", "9400ZZMASTP3", "9400ZZMASTP1"}

		assert.Equal(t, expectedStopsInAreaMap, stopsInAreaMap)
	})

	t.Run(`Given NaPTAN CSV data cannot be retrieved
When FetchStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		stopsInAreaFilename := "StopsInArea.csv"
		stopAreaCodeColumnIndex := 0
		atcoCodeColumnIndex := 1

		zipFileFetcher := mock_http.NewMockZipFileFetcher(ctrl)
		zipFileFetcherErr := errors.New("FUBAR")
		zipFileFetcher.EXPECT().FetchZipFile(ctx).Return(nil, zipFileFetcherErr)

		extractor := mock_compression.NewMockExtractor(ctrl)

		naptanCsv := naptan.NewCSV(logger, zipFileFetcher, extractor, stopsInAreaFilename, stopAreaCodeColumnIndex, atcoCodeColumnIndex)

		// When
		stopsInAreaMap, err := naptanCsv.FetchStopsInArea(ctx)

		// Then
		assert.Nil(t, stopsInAreaMap)
		assert.NotNil(t, err)
		assert.Equal(t, zipFileFetcherErr, err)
	})

	t.Run(`Given NaPTAN CSV data cannot be extracted
When FetchStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		zipData := []byte("arbitrary data")

		stopsInAreaFilename := "StopsInArea.csv"
		stopAreaCodeColumnIndex := 0
		atcoCodeColumnIndex := 1

		zipFileFetcher := mock_http.NewMockZipFileFetcher(ctrl)
		zipFileFetcher.EXPECT().FetchZipFile(ctx).Return(givenArbitraryZipData(t, zipData), nil)

		extractor := mock_compression.NewMockExtractor(ctrl)
		extractorErr := errors.New("FUBAR")
		extractor.EXPECT().ExtractFile(zipData, stopsInAreaFilename).Return(nil, extractorErr)

		naptanCsv := naptan.NewCSV(logger, zipFileFetcher, extractor, stopsInAreaFilename, stopAreaCodeColumnIndex, atcoCodeColumnIndex)

		// When
		stopsInAreaMap, err := naptanCsv.FetchStopsInArea(ctx)

		// Then
		assert.Nil(t, stopsInAreaMap)
		assert.NotNil(t, err)
		assert.Equal(t, extractorErr, err)
	})

	t.Run(`Given NaPTAN CSV data is empty
When FetchStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		zipData := []byte("arbitrary data")

		stopsInAreaFilename := "StopsInArea.csv"
		stopAreaCodeColumnIndex := 0
		atcoCodeColumnIndex := 1

		zipFileFetcher := mock_http.NewMockZipFileFetcher(ctrl)
		zipFileFetcher.EXPECT().FetchZipFile(ctx).Return(givenArbitraryZipData(t, zipData), nil)

		extractor := mock_compression.NewMockExtractor(ctrl)
		extractor.EXPECT().ExtractFile(zipData, stopsInAreaFilename).Return(givenEmptyCsvData(t), nil)

		naptanCsv := naptan.NewCSV(logger, zipFileFetcher, extractor, stopsInAreaFilename, stopAreaCodeColumnIndex, atcoCodeColumnIndex)

		// When
		stopsInAreaMap, err := naptanCsv.FetchStopsInArea(ctx)

		// Then
		assert.Nil(t, stopsInAreaMap)
		assert.NotNil(t, err)
		assert.Equal(t, io.EOF, err)
	})

	t.Run(`Given NaPTAN CSV data is corrupt
When FetchStopsInArea is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		logger := mockLogger(t)

		zipData := []byte("arbitrary data")

		stopsInAreaFilename := "StopsInArea.csv"
		stopAreaCodeColumnIndex := 0
		atcoCodeColumnIndex := 1

		zipFileFetcher := mock_http.NewMockZipFileFetcher(ctrl)
		zipFileFetcher.EXPECT().FetchZipFile(ctx).Return(givenArbitraryZipData(t, zipData), nil)

		extractor := mock_compression.NewMockExtractor(ctrl)
		extractor.EXPECT().ExtractFile(zipData, stopsInAreaFilename).Return(givenCorruptCsvData(t), nil)

		naptanCsv := naptan.NewCSV(logger, zipFileFetcher, extractor, stopsInAreaFilename, stopAreaCodeColumnIndex, atcoCodeColumnIndex)

		// When
		stopsInAreaMap, err := naptanCsv.FetchStopsInArea(ctx)

		// Then
		assert.Nil(t, stopsInAreaMap)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "record on line 2: wrong number of fields")
	})
}
