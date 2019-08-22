package main

import (
	"fmt"
	"net/http"

	"github.com/metrics/listener"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	_ "net/http/pprof"

	"github.com/metrics/db"
	"github.com/metrics/writer"
)

func main() {
	var (
		dbFilePath = pflag.String("file_path", "/tmp/db.txt", "Db file path. It will be created in case of absence")

		listenerHost = pflag.String("ip", "0.0.0.0", "Listener host")
		listenerPort = pflag.Int("port", 9000, "Listener port")

		metricsFlushPeriod = pflag.Int("flush_period", 5, "Metrics flush period in seconds")
	)
	pflag.Parse()

	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf(":%d", 9001), nil))
	}()

	dbConn, err := db.New(*dbFilePath)
	if err != nil {
		log.Fatalf("Cannot initialize db connection. Err: %s", err)
	}

	metricsWriter := writer.NewMetricsWriter(dbConn, *metricsFlushPeriod)
	go metricsWriter.Run()

	metricsListener := listener.NewListener(*listenerHost, *listenerPort, metricsWriter)
	if err := metricsListener.Run(); err != nil {
		log.Fatalf("Listener has been stopped with error. Err: %s", err)
	}
}
