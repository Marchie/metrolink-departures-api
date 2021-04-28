package naptan

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Repository struct {
	logger     *zap.Logger
	httpClient *http.Client
	url        string
}

func NewRepository(logger *zap.Logger, httpClient *http.Client, url string) *Repository {
	return &Repository{
		logger:     logger,
		httpClient: httpClient,
		url:        url,
	}
}

func (r *Repository) FetchZipFile(ctx context.Context) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.url, nil)
	if err != nil {
		return nil, err
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from %s: %s", r.url, res.Status)
	}

	return res.Body, nil
}
