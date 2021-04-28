package apigw_test

import (
	"bytes"
	"context"
	mock_core "github.com/Marchie/tf-experiment/lambda/internal/mocks/core"
	"github.com/Marchie/tf-experiment/lambda/internal/transport/apigw"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"io/ioutil"
	"net/http"
	"testing"
)

func thenExpHeaders(t *testing.T) map[string]string {
	t.Helper()

	headers := make(map[string]string)
	headers["Content-type"] = "application/json"

	return headers
}

func TestMetrolinkDeparturesAwsApiGateway_Handler(t *testing.T) {
	t.Run(`Given a configured Metrolink Departures AWS API Gateway
When Handler is called with a StopAreaCode in the path parameter
Then Metrolink departures for the StopAreaCode are returned in the response`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		stopAreaCode := "940GZZMASTP"

		zapCore, observedLogs := observer.New(zap.ErrorLevel)
		logger := zap.New(zapCore)

		apiData := `{
	"stopAreaCode": "940GZZMASTP",
	"departures": [
		{
			"sequence": 0,
			"destination": "Altrincham",
			"wait": "2",
			"status": "Due",
			"carriages": "Double",
			"platform": "A"
		},
		{
			"sequence": 1,
			"destination": "Victoria",
			"wait": "4",
			"status": "Due",
			"carriages": "Single",
			"platform": "D"
		}
	],
	"lastUpdated": "2021-03-24T21:26:52Z"
}`

		stopAreaCodePathParameter := "stopAreaCode"
		pathParameters := make(map[string]string)
		pathParameters[stopAreaCodePathParameter] = stopAreaCode

		metrolinkDeparturesJsonApi := mock_core.NewMockStopAreaDeparturesJsoner(ctrl)
		metrolinkDeparturesJsonApi.EXPECT().Json(ctx, stopAreaCode).Return(ioutil.NopCloser(bytes.NewBufferString(apiData)), http.StatusOK, nil)

		metrolinkDeparturesAwsApiGateway := apigw.NewMetrolinkDeparturesAwsApiGateway(logger, metrolinkDeparturesJsonApi, stopAreaCodePathParameter)

		apiGatewayProxyRequest := events.APIGatewayProxyRequest{
			PathParameters: pathParameters,
		}

		// When
		apiGatewayProxyResponse, err := metrolinkDeparturesAwsApiGateway.Handler(ctx, apiGatewayProxyRequest)

		// Then
		assert.Nil(t, err)

		expApiGatewayProxyResponse := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    thenExpHeaders(t),
			Body:       apiData,
		}

		assert.EqualValues(t, expApiGatewayProxyResponse, apiGatewayProxyResponse)

		assert.Equal(t, 0, observedLogs.Len())
	})

	t.Run(`Given a configured Metrolink Departures AWS API Gateway
When Handler is called without a StopAreaCode in the path parameter
Then an error response is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		zapCore, observedLogs := observer.New(zap.ErrorLevel)
		logger := zap.New(zapCore)

		stopAreaCodePathParameter := "stopAreaCode"
		pathParameters := make(map[string]string)

		metrolinkDeparturesJsonApi := mock_core.NewMockStopAreaDeparturesJsoner(ctrl)

		metrolinkDeparturesAwsApiGateway := apigw.NewMetrolinkDeparturesAwsApiGateway(logger, metrolinkDeparturesJsonApi, stopAreaCodePathParameter)

		apiGatewayProxyRequest := events.APIGatewayProxyRequest{
			PathParameters: pathParameters,
		}

		// When
		apiGatewayProxyResponse, err := metrolinkDeparturesAwsApiGateway.Handler(ctx, apiGatewayProxyRequest)

		// Then
		assert.Nil(t, err)

		expApiGatewayProxyResponse := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    thenExpHeaders(t),
			Body: `{
	"error": "no StopAreaCode in request path parameters"
}`,
		}

		assert.EqualValues(t, expApiGatewayProxyResponse, apiGatewayProxyResponse)

		assert.Equal(t, 1, observedLogs.Len())
		assert.Equal(t, zapcore.ErrorLevel, observedLogs.All()[0].Level)
		assert.Equal(t, "no StopAreaCode or AtcoCode in request path parameters", observedLogs.All()[0].Message)
	})

	t.Run(`Given an error occurs retrieving Metrolink Departures API data
When Handler is called with a StopAreaCode in the path parameter
Then an error response is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		stopAreaCode := "940GZZMASTP"

		zapCore, observedLogs := observer.New(zap.ErrorLevel)
		logger := zap.New(zapCore)

		stopAreaCodePathParameter := "stopAreaCode"
		pathParameters := make(map[string]string)
		pathParameters[stopAreaCodePathParameter] = stopAreaCode

		metrolinkDeparturesJsonApi := mock_core.NewMockStopAreaDeparturesJsoner(ctrl)
		metrolinkDeparturesJsonApiErr := errors.New("FUBAR")
		metrolinkDeparturesJsonApi.EXPECT().Json(ctx, stopAreaCode).Return(nil, http.StatusInternalServerError, metrolinkDeparturesJsonApiErr)

		metrolinkDeparturesAwsApiGateway := apigw.NewMetrolinkDeparturesAwsApiGateway(logger, metrolinkDeparturesJsonApi, stopAreaCodePathParameter)

		apiGatewayProxyRequest := events.APIGatewayProxyRequest{
			PathParameters: pathParameters,
		}

		// When
		apiGatewayProxyResponse, err := metrolinkDeparturesAwsApiGateway.Handler(ctx, apiGatewayProxyRequest)

		// Then
		assert.Nil(t, err)

		expApiGatewayProxyResponse := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    thenExpHeaders(t),
			Body: `{
	"error": "internal server error"
}`,
		}

		assert.EqualValues(t, expApiGatewayProxyResponse, apiGatewayProxyResponse)

		assert.Equal(t, 1, observedLogs.Len())
		assert.Equal(t, zapcore.ErrorLevel, observedLogs.All()[0].Level)
		assert.Equal(t, "error with Metrolink Departures API JSON response", observedLogs.All()[0].Message)
		assert.Equal(t, "stopAreaCode", observedLogs.All()[0].Context[0].Key)
		assert.Equal(t, stopAreaCode, observedLogs.All()[0].Context[0].String)
		assert.Equal(t, "error", observedLogs.All()[0].Context[1].Key)
		assert.Equal(t, metrolinkDeparturesJsonApiErr, observedLogs.All()[0].Context[1].Interface)
	})

	t.Run(`Given an empty response from the Metrolink Departures API
When Handler is called with a StopAreaCode in the path parameter
Then an error response is returned`, func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		stopAreaCode := "940GZZMASTP"

		zapCore, observedLogs := observer.New(zap.ErrorLevel)
		logger := zap.New(zapCore)

		stopAreaCodePathParameter := "stopAreaCode"
		pathParameters := make(map[string]string)
		pathParameters[stopAreaCodePathParameter] = stopAreaCode

		metrolinkDeparturesJsonApi := mock_core.NewMockStopAreaDeparturesJsoner(ctrl)
		metrolinkDeparturesJsonApi.EXPECT().Json(ctx, stopAreaCode).Return(ioutil.NopCloser(bytes.NewBufferString("")), http.StatusOK, nil)

		metrolinkDeparturesAwsApiGateway := apigw.NewMetrolinkDeparturesAwsApiGateway(logger, metrolinkDeparturesJsonApi, stopAreaCodePathParameter)

		apiGatewayProxyRequest := events.APIGatewayProxyRequest{
			PathParameters: pathParameters,
		}

		// When
		apiGatewayProxyResponse, err := metrolinkDeparturesAwsApiGateway.Handler(ctx, apiGatewayProxyRequest)

		// Then
		assert.Nil(t, err)

		expApiGatewayProxyResponse := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    thenExpHeaders(t),
			Body: `{
	"error": "no data for StopAreaCode 940GZZMASTP"
}`,
		}

		assert.EqualValues(t, expApiGatewayProxyResponse, apiGatewayProxyResponse)

		assert.Equal(t, 0, observedLogs.Len())
	})
}
