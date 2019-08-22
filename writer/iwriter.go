package writer

import (
	"context"

	"github.com/metrics/models"
)

type IMetricWriter interface {
	Write(ctx context.Context, metric *models.Metric) error
	Run()
}
