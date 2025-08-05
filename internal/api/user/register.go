package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/internal/metric"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

// PostUserRegister - регистрация пользователя
func (i *Implementation) PostUserRegister(w http.ResponseWriter, r *http.Request) {
	// Объявляем структуру для хранения входных данных
	var info *api.PostUserRegisterJSONBody

	metric.IncRequestCounter()
	timeStart := time.Now()

	// Парсим тело запроса в структуру info
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PostUserRegister")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	userInfo := converter.ToUserInfoFromApi(info)
	id, err := i.userService.Register(context.Background(), userInfo)
	diffTime := time.Since(timeStart)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PostUserRegister")
		metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("inserted user with id: %d", id)

	response := struct {
		UserId string `json:"userId"`
	}{
		UserId: id,
	}

	// Отправляем объект userObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PostUserRegister")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PostUserRegister")
	metric.HistogramResponseTimeObserve("success", diffTime.Seconds())
}
