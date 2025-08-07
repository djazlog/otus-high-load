package api

import (
	"context"
	"encoding/json"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/internal/metric"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

func (i *Implementation) GetPostFeed(w http.ResponseWriter, r *http.Request, params api.GetPostFeedParams) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	userObj, err := i.postService.Get(context.Background(), params.Limit, params.Offset)
	diffTime := time.Since(timeStart)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "GetPostFeed")
		metric.HistogramResponseTimeObserve("GetPostFeed", diffTime.Seconds())
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToPostFromService(userObj)
	// Отправляем объект userObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetPostFeed")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetPostFeed")
	metric.HistogramResponseTimeObserve("GetPostFeed", diffTime.Seconds())
}
