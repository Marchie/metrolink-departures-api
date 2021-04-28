package compression_test

import (
	"archive/zip"
	"bytes"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/compression"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"io"
	"io/ioutil"
	"testing"
)

func givenZippedData(t *testing.T) []byte {
	t.Helper()

	zipArchive := new(bytes.Buffer)
	filesToZip := make(map[string][]byte)
	filesToZip["file1.txt"] = []byte("The quick brown fox jumped over the lazy dog.")
	filesToZip["file2.txt"] = []byte("The slow grey cat rode on the fox's back.")

	zipWriter := zip.NewWriter(zipArchive)

	for fileName, data := range filesToZip {
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			t.Fatal(err)
		}

		_, err = fileWriter.Write(data)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}

	zipData, err := ioutil.ReadAll(zipArchive)
	if err != nil {
		t.Fatal(err)
	}

	return zipData
}

func mockLogger(t *testing.T) *zap.Logger {
	t.Helper()

	zapCore, _ := observer.New(zapcore.DebugLevel)
	return zap.New(zapCore)
}

func readData(t *testing.T, rc io.ReadCloser) []byte {
	t.Helper()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func TestFileExtractor_ExtractFile(t *testing.T) {
	t.Run(`Given zipped data containing a file
When ExtractFile is called with the file name
Then an io.ReadCloser containing the file contents is returned`, func(t *testing.T) {
		// Given
		logger := mockLogger(t)

		zipData := givenZippedData(t)

		extractor := compression.NewZipFileExtractor(logger)

		// When
		rc, err := extractor.ExtractFile(zipData, "file1.txt")

		// Then
		assert.NotNil(t, rc)
		assert.Nil(t, err)
		assert.Equal(t, []byte("The quick brown fox jumped over the lazy dog."), readData(t, rc))
	})

	t.Run(`Given zipped data does not contain a file
When ExtractFile is called with the file name
Then an error is returned`, func(t *testing.T) {
		// Given
		logger := mockLogger(t)

		zipData := givenZippedData(t)

		extractor := compression.NewZipFileExtractor(logger)

		// When
		rc, err := extractor.ExtractFile(zipData, "file3.txt")

		// Then
		assert.Nil(t, rc)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "file 'file3.txt' not found in zip archive")
	})

	t.Run(`Given a non-zip data
When ExtractFile is called
Then an error is returned`, func(t *testing.T) {
		// Given
		logger := mockLogger(t)

		zipData := []byte("not zip data")

		extractor := compression.NewZipFileExtractor(logger)

		// When
		rc, err := extractor.ExtractFile(zipData, "file1.txt")

		// Then
		assert.Nil(t, rc)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "zip: not a valid zip file")
	})
}
