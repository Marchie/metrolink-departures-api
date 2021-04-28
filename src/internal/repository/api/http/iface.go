package http

import (
	"context"
	"io"
)

type ZipFileFetcher interface {
	FetchZipFile(ctx context.Context) (io.ReadCloser, error)
}
