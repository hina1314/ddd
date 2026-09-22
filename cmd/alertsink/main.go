package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type webhookPayload struct {
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	Alerts            []webhookAlert    `json:"alerts"`
}

type webhookAlert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
	Fingerprint string            `json:"fingerprint"`
}

func main() {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).
			With("service", "alertsink"),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /alerts", receiveAlerts)

	server := &http.Server{
		Addr:              "127.0.0.1:5001",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	slog.Info("alertsink_started", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("alertsink_stopped", "error", err)
		os.Exit(1)
	}
}

func receiveAlerts(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var payload webhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Warn("invalid_alert_payload", "error", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	slog.Info(
		"alert_received",
		"status", payload.Status,
		"receiver", payload.Receiver,
		"group_labels", payload.GroupLabels,
		"common_labels", payload.CommonLabels,
		"common_annotations", payload.CommonAnnotations,
		"alert_count", len(payload.Alerts),
		"alerts", payload.Alerts,
	)

	w.WriteHeader(http.StatusNoContent)
}
