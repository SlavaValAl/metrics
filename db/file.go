package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/metrics/models"
	"github.com/pkg/errors"
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
	file, err := os.OpenFile(
		f.filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return errors.WithMessage(err, "cannot open db connection")
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Cannot close db connection. Err: %s", err)
		}
	}()

	writer := bufio.NewWriter(file)
	for _, metric := range metrics {
		metricJson, err := json.Marshal(metric)
		if err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot convert metric %v to json format", metric))
		}

		if _, err := writer.Write(metricJson); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot save metric %v to db", metric))
		}
		if _, err := writer.WriteString("\n"); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("cannot save metric %v to db", metric))
		}
	}

	if err := writer.Flush(); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot write metrics to db"))
	}

	return nil
}
