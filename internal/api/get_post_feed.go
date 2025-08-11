package api

import (
	"encoding/json"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/internal/metric"
	"otus-project/internal/utils"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

// GetPostFeed GET /post/feed
func (i *Implementation) GetPostFeed(w http.ResponseWriter, r *http.Request, params api.GetPostFeedParams) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	defer func() {
		diffTime := time.Since(timeStart)
		metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetPostFeed")
		metric.HistogramResponseTimeObserve("GetPostFeed", diffTime.Seconds())
	}()

	userId, err := utils.GetUserFromToken(r)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	postsObj, err := i.postService.Feed(r.Context(), *userId, params.Offset, params.Limit)

	if err != nil {

		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToPostsFromService(postsObj)
	// Отправляем объект userObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
