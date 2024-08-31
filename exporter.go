package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

type metrics []map[string]string

func gatherMetrics() metrics {
	hostname, _ := os.Hostname()

	return metrics{{
		"hostname": hostname,
		"name":     "cpu-num",
		"value":    strconv.Itoa(runtime.NumCPU()),
	}}
}

func SendMetricsToServer(interval int, address string) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		buffer := new(bytes.Buffer)

		err := json.NewEncoder(buffer).Encode(gatherMetrics())
		if err != nil {
			slog.Error("exporter", "error", err.Error())
			continue
		}

		response, err := http.Post(fmt.Sprintf("http://%s/push", address), "application/json", buffer)
		if err != nil {
			slog.Error("exporter", "error", err.Error())
		} else {
			slog.Info("exporter", "message", fmt.Sprintf("Sending metrics to %s", address))
		}
		response.Body.Close()
	}
}
