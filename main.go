package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	isExporter := flag.Bool("exporter", false, "Runs in exporter mode")
	isServer := flag.Bool("server", false, "Runs in server mode")
	address := flag.String("address", "127.0.0.1:3333", "Address to listen (server) or where to push (exporter)")
	interval := flag.Int("interval", 5, "Interval of sending metrics to server")

	flag.Parse()

	if *isExporter {
		if *isServer {
			go SendMetricsToServer(*interval, *address)
		} else {
			SendMetricsToServer(*interval, *address)
			return
		}
	}

	if *isServer {
		db, err := GetDB()
		if err != nil {
			slog.Error("server", "error", err.Error())
			return
		}

		http.HandleFunc("/metrics", GetMetrics(db))
		http.HandleFunc("/push", PushMetric(db))
		slog.Info("server", "message", fmt.Sprintf("Listening %s", *address))

		err = http.ListenAndServe(*address, nil)
		if err != nil {
			slog.Error("server", "error", err.Error())
		}

		return
	}

	flag.PrintDefaults()
}
