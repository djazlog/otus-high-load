package api

import (
	"encoding/json"
	"net/http"
	"otus-project/internal/metric"
	"strconv"
	"time"
)

// GetDialogStats GET /dialog/stats
func (i *Implementation) GetDialogStats(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	// Получаем статистику диалогов
	totalDialogs, activeDialogs, err := i.dialogService.GetDialogStats(r.Context())
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetDialogStats")
		metric.HistogramResponseTimeObserve("GetDialogStatsError", diffTime.Seconds())
		http.Error(w, "Failed to get dialog stats", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"total_dialogs":    totalDialogs,
		"active_dialogs":   activeDialogs,
		"response_time_ms": diffTime.Milliseconds(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetDialogStats")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetDialogStats")
	metric.HistogramResponseTimeObserve("GetDialogStats", diffTime.Seconds())
}

// GetDialogCount GET /dialog/count
func (i *Implementation) GetDialogCount(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	// Получаем количество диалогов
	count, err := i.dialogService.GetDialogCount(r.Context())
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetDialogCount")
		metric.HistogramResponseTimeObserve("GetDialogCountError", diffTime.Seconds())
		http.Error(w, "Failed to get dialog count", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"dialog_count":     count,
		"response_time_ms": diffTime.Milliseconds(),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetDialogCount")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetDialogCount")
	metric.HistogramResponseTimeObserve("GetDialogCount", diffTime.Seconds())
}
