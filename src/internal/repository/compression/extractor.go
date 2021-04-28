package compression

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
)

type ZipFileExtractor struct {
	logger *zap.Logger
}

func NewZipFileExtractor(logger *zap.Logger) *ZipFileExtractor {
	return &ZipFileExtractor{
		logger: logger,
	}
}

func (c *ZipFileExtractor) ExtractFile(zipData []byte, filename string) (io.ReadCloser, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	return c.getFileFromZipArchive(zipReader, filename)
}

func (c *ZipFileExtractor) getFileFromZipArchive(zipReader *zip.Reader, filename string) (io.ReadCloser, error) {
	for _, file := range zipReader.File {
		if file.Name == filename {
			return c.getReadCloserOfZippedFile(file)
		}
	}

	return nil, fmt.Errorf("file '%s' not found in zip archive", filename)
}

func (c *ZipFileExtractor) getReadCloserOfZippedFile(file *zip.File) (io.ReadCloser, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}

	return rc, nil
}
