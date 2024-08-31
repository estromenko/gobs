package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func GetMetrics(db *sql.DB) http.HandlerFunc {
	type metric struct {
		Date     string `json:"date"`
		Hostname string `json:"hostname"`
		Name     string `json:"name"`
		Value    string `json:"value"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []metric

		rows, err := db.Query("SELECT date, hostname, name, value FROM metrics")
		if err != nil {
			slog.Error(err.Error())
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		for rows.Next() {
			var metric metric
			err = rows.Scan(&metric.Date, &metric.Hostname, &metric.Name, &metric.Value)
			if err != nil {
				slog.Error(err.Error())
				_ = json.NewEncoder(w).Encode(map[string]string{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}
			metrics = append(metrics, metric)
		}

		_ = json.NewEncoder(w).Encode(metrics)
	}
}

func PushMetric(db *sql.DB) http.HandlerFunc {
	type request []struct {
		Hostname string `json:"hostname"`
		Name     string `json:"name"`
		Value    string `json:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Error("server", "error", err.Error())
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		var insertValues []string
		for _, metric := range req {
			insertValues = append(insertValues, fmt.Sprintf(`("%s", "%s", "%s")`, metric.Hostname, metric.Name, metric.Value))
		}

		query := "INSERT INTO metrics (hostname, name, value) VALUES " + strings.Join(insertValues, ",")
		_, err = db.Exec(query)
		if err != nil {
			slog.Error("server", "error", err.Error())
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		slog.Info("server", "message", "Metrics received")

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
			"error":  nil,
		})
	}
}
