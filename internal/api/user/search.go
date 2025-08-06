package user

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

func (i *Implementation) GetUserSearch(w http.ResponseWriter, r *http.Request, params api.GetUserSearchParams) {
	w.Header().Set("Content-Type", "application/json")

	metric.IncRequestCounter()
	timeStart := time.Now()
	filter := converter.ToUserFilterFromApi(&params)

	usersObj, err := i.userService.Search(context.Background(), filter)
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "GetUserSearch")
		metric.HistogramResponseTimeObserve("GetUserSearchError", diffTime.Seconds())
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToUsersFromService(usersObj)

	// Отправляем объект usersObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetUserSearchError")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetUserSearch")
	metric.HistogramResponseTimeObserve("GetUserSearch", diffTime.Seconds())
}
