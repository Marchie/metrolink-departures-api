package sqs

import (
	"context"
	"github.com/Marchie/tf-experiment/lambda/internal/core"
	"go.uber.org/zap"
)

type MetrolinkDeparturesDataLoader struct {
	logger                    *zap.Logger
	metrolinkDeparturesLoader core.MetrolinkDeparturesLoader
}

func NewMetrolinkDeparturesDataLoader(logger *zap.Logger, metrolinkDeparturesLoader core.MetrolinkDeparturesLoader) *MetrolinkDeparturesDataLoader {
	return &MetrolinkDeparturesDataLoader{
		logger:                    logger,
		metrolinkDeparturesLoader: metrolinkDeparturesLoader,
	}
}

func (m *MetrolinkDeparturesDataLoader) Handler(ctx context.Context) error {
	return m.metrolinkDeparturesLoader.Load(ctx)
}
