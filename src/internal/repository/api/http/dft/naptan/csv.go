package naptan

import (
	"context"
	"encoding/csv"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/api/http"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/compression"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

type CSV struct {
	logger                  *zap.Logger
	zipFileFetcher          http.ZipFileFetcher
	extractor               compression.Extractor
	stopsInAreaFilename     string
	stopAreaCodeColumnIndex int
	atcoCodeColumnIndex     int
}

func NewCSV(logger *zap.Logger, zipFileFetcher http.ZipFileFetcher, extractor compression.Extractor, stopsInAreaFilename string, stopAreaCodeColumnIndex int, atcoCodeColumnIndex int) *CSV {
	return &CSV{
		logger:                  logger,
		zipFileFetcher:          zipFileFetcher,
		extractor:               extractor,
		stopsInAreaFilename:     stopsInAreaFilename,
		stopAreaCodeColumnIndex: stopAreaCodeColumnIndex,
		atcoCodeColumnIndex:     atcoCodeColumnIndex,
	}
}

func (c *CSV) FetchStopsInArea(ctx context.Context) (map[string][]string, error) {
	zipFile, err := c.zipFileFetcher.FetchZipFile(ctx)
	if err != nil {
		return nil, err
	}

	zipData, err := ioutil.ReadAll(zipFile)
	if err != nil {
		return nil, err
	}

	stopsInAreaReadCloser, err := c.extractor.ExtractFile(zipData, c.stopsInAreaFilename)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(stopsInAreaReadCloser)
	if err := c.skipHeaderRow(csvReader); err != nil {
		return nil, err
	}

	stopsInArea := make(map[string][]string)

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		stopAreaCode, atcoCode := row[c.stopAreaCodeColumnIndex], row[c.atcoCodeColumnIndex]

		stopsInArea[stopAreaCode] = append(stopsInArea[stopAreaCode], atcoCode)
	}

	return stopsInArea, nil
}

func (*CSV) skipHeaderRow(csvReader *csv.Reader) error {
	_, err := csvReader.Read()
	return err
}
