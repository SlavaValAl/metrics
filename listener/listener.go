package listener

import (
	"encoding/json"
	"fmt"
	"github.com/metrics/writer"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/metrics/models"
)

type Listener struct {
	host       string
	port       int
	reqTimeout int

	metricWriter writer.IMetricWriter
}

func NewListener(host string, port int, metricWriter writer.IMetricWriter) *Listener {
	return &Listener{
		host:         host,
		port:         port,
		metricWriter: metricWriter,
	}
}

func (l *Listener) Run() error {
	http.HandleFunc("/ga", l.newMetric)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", l.host, l.port), nil)
}

func (l *Listener) newMetric(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer func () {
		if err := r.Body.Close(); err != nil {
			log.Errorf("Cannot close response body. Err: %s", err)
		}
	}()

	var metric models.Metric

	err := decoder.Decode(&metric)
	if err != nil {
		log.Errorf("Cannot decode metric. Err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		if _, sndErr := w.Write([]byte(fmt.Sprintf("Cannot decode metric. Err: %s", err))); sndErr != nil {
			log.Errorf("Cannot write response message. Err: %s", sndErr)
		}
		return
	}

	if err := l.metricWriter.Write(r.Context(), &metric); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, sndErr := w.Write([]byte(fmt.Sprintf("Metric can not be saved due to tool error. Err: %s", err))); sndErr != nil {
			log.Errorf("Cannot write response message. Err: %s", sndErr)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, sndErr := w.Write([]byte("Metric has been successfully saved")); sndErr != nil {
		log.Errorf("Cannot write response message. Err: %s", sndErr)
	}
}
