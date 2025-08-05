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

// GetUserGetId - получение пользователя по id
func (i *Implementation) GetUserGetId(w http.ResponseWriter, r *http.Request, id api.UserId) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	userObj, err := i.userService.Get(context.Background(), id)
	diffTime := time.Since(timeStart)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "GetUserById")
		metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToUserFromService(userObj)
	// Отправляем объект userObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetUserById")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetUserById")
}
