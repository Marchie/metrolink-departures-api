package apigw

import (
	"context"
	"fmt"
	"github.com/Marchie/tf-experiment/lambda/internal/core"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type MetrolinkDeparturesAwsApiGateway struct {
	logger                              *zap.Logger
	stopAreaDeparturesJsoner            core.StopAreaDeparturesJsoner
	stopAreaCodeOrAtcoCodePathParameter string
}

func NewMetrolinkDeparturesAwsApiGateway(logger *zap.Logger, jsoner core.StopAreaDeparturesJsoner, stopAreaCodeOrAtcoCodePathParameter string) *MetrolinkDeparturesAwsApiGateway {
	return &MetrolinkDeparturesAwsApiGateway{
		logger:                              logger,
		stopAreaDeparturesJsoner:            jsoner,
		stopAreaCodeOrAtcoCodePathParameter: stopAreaCodeOrAtcoCodePathParameter,
	}
}

func (h *MetrolinkDeparturesAwsApiGateway) Handler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"

	stopAreaCodeOrAtcoCode := event.PathParameters[h.stopAreaCodeOrAtcoCodePathParameter]

	if stopAreaCodeOrAtcoCode == "" {
		h.logger.Error("no StopAreaCode or AtcoCode in request path parameters")

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body: `{
	"error": "no StopAreaCode in request path parameters"
}`,
		}, nil
	}

	departures, statusCode, err := h.stopAreaDeparturesJsoner.Json(ctx, stopAreaCodeOrAtcoCode)
	if err != nil {
		h.logger.Error("error with Metrolink Departures API JSON response", zap.String("stopAreaCode", stopAreaCodeOrAtcoCode), zap.Error(err))

		return &events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Headers:    headers,
			Body: `{
	"error": "internal server error"
}`,
		}, nil
	}

	buf := new(strings.Builder)
	written, err := io.Copy(buf, departures)
	if err != nil {
		h.logger.Error("error reading Metrolink Departures API JSON response", zap.String("stopAreaCode", stopAreaCodeOrAtcoCode), zap.Error(err))

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body: `{
	"error": "internal server error"
}`,
		}, nil
	}

	if written == 0 {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    headers,
			Body: fmt.Sprintf(`{
	"error": "no data for StopAreaCode %s"
}`, stopAreaCodeOrAtcoCode),
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       buf.String(),
	}, nil
}
