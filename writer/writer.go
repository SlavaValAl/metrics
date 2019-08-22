package writer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/metrics/db"
	"github.com/metrics/models"
)

type MetricWriter struct {
	writePeriod int
	db          db.Db

	metricsChan chan *models.Metric

	metrics []*models.Metric
}

func NewMetricsWriter(db db.Db, writePeriod int) IMetricWriter {
	return &MetricWriter{
		db:          db,
		writePeriod: writePeriod,
		metricsChan: make(chan *models.Metric),
		metrics:     make([]*models.Metric, 0),
	}
}

func (w *MetricWriter) Write(ctx context.Context, metric *models.Metric) error {
	select {
	case w.metricsChan <- metric:
	case <-ctx.Done():
		return errors.New("Cancelled on timeout")
	}
	return nil
}

func (w *MetricWriter) Run() {
	ticker := time.NewTicker(time.Duration(w.writePeriod) * time.Second)

	for {
		select {
		case <-ticker.C:
			if len(w.metrics) == 0 {
				continue
			}

			w.flush(w.metrics)
			w.metrics = make([]*models.Metric, 0)
		case metric := <-w.metricsChan:
			if metric.IsEmpty() {
				log.Info("Empty metric will be ignored")
				continue
			}
			w.metrics = append(w.metrics, metric)
		}
	}
}

func (w *MetricWriter) flush(metrics []*models.Metric) {
	if err := w.db.SaveMetrics(metrics); err != nil {
		log.Errorf("Cannot save metrics. Err: %s", err)
	}
}
