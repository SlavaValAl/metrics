package db

import "github.com/metrics/models"

type Db interface {
	SaveMetrics(metrics []*models.Metric) error
}
