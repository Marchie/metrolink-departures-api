package naptan_test

import (
	"context"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/repository/api/http/dft/naptan"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockHttpClient(t *testing.T) *http.Client {
	t.Helper()

	return &http.Client{}
}

func givenWorkingNaptanEndpoint(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("arbitrary data"))
	}))
}

func givenNaptanEndpointReturns404(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("page not found"))
	}))
}

func readData(t *testing.T, rc io.ReadCloser) []byte {
	t.Helper()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func TestRepository_FetchZipFile(t *testing.T) {
	t.Run(`Given a HTTP endpoint containing NaPTAN data
When FetchZipFile is called
Then the file is returned`, func(t *testing.T) {
		// Given
		ctx := context.Background()

		naptanServer := givenWorkingNaptanEndpoint(t)
		defer naptanServer.Close()

		logger := mockLogger(t)

		repository := naptan.NewRepository(logger, mockHttpClient(t), naptanServer.URL)

		// When
		rc, err := repository.FetchZipFile(ctx)

		// Then
		assert.NotNil(t, rc)
		assert.Nil(t, err)
		assert.Equal(t, []byte("arbitrary data"), readData(t, rc))
	})

	t.Run(`Given a HTTP endpoint cannot be reached
When FetchZipFile is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctx := context.Background()

		logger := mockLogger(t)

		repository := naptan.NewRepository(logger, mockHttpClient(t), "http://invalid")

		// When
		rc, err := repository.FetchZipFile(ctx)

		// Then
		assert.Nil(t, rc)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "Get \"http://invalid\": dial tcp: lookup invalid: no such host")
	})

	t.Run(`Given a HTTP endpoint returns an error
When FetchZipFile is called
Then an error is returned`, func(t *testing.T) {
		// Given
		ctx := context.Background()

		naptanServer := givenNaptanEndpointReturns404(t)
		defer naptanServer.Close()

		logger := mockLogger(t)

		repository := naptan.NewRepository(logger, mockHttpClient(t), naptanServer.URL)

		// When
		rc, err := repository.FetchZipFile(ctx)

		// Then
		assert.Nil(t, rc)
		assert.NotNil(t, err)
		assert.EqualError(t, err, fmt.Sprintf("error response from %s: %s", naptanServer.URL, "404 Not Found"))
	})
}
