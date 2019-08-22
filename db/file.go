package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/metrics/models"
)

type FileDb struct {
	filePath string
}

func New(filePath string) (Db, error) {
	db := &FileDb{filePath: filePath}
	if err := db.checkConnection(); err != nil {
		return nil, err
	}
	return db, nil
}

func (f *FileDb) checkConnection() error {
	file, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func (f *FileDb) SaveMetrics(metrics []*models.Metric) error {
	var metricsBuffer bytes.Buffer
	for _, metric := range metrics {
		metricJson, err := json.Marshal(metric)
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot convert metric %v to json format", metric))
		}

		if _, err := metricsBuffer.Write(metricJson); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot save metric %v to db", metric))
		}
		if _, err := metricsBuffer.WriteString("\n"); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot save metric %v to db", metric))
		}
	}
	log.Debugf("Metrics to save: %s", metricsBuffer.String())

	file, err := os.OpenFile(
		f.filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return errors.WithMessage(err, "cannot open db connection")
	}

	if _, err := file.Write(metricsBuffer.Bytes()); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot save metrics %v to db", metrics))
	}

	if err := file.Close(); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot close db connection after the metrics writing"))
	}

	return nil
}
