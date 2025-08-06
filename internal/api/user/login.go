package user

import (
	"context"
	"encoding/json"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/internal/metric"
	"otus-project/internal/model"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

// PostLogin - обработчик POST запроса на /login
func (i *Implementation) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Объявляем структуру для хранения входных данных
	var info *api.PostLoginJSONRequestBody

	metric.IncRequestCounter()
	timeStart := time.Now()
	// Парсим тело запроса в структуру info
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PostLogin")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	loginDto := &model.LoginDto{Id: *info.Id, Password: *info.Password}
	/// TODO: Валидация

	loginObj, err := i.userService.Login(context.Background(), loginDto)
	diffTime := time.Since(timeStart)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "PostLogin")
		metric.HistogramResponseTimeObserve("PostLoginError", diffTime.Seconds())
		http.Error(w, "Login failed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	tokenResponse := converter.ToTokenResponse(loginObj)
	// Отправляем объект userObj в формате JSON
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PostLogin")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PostLogin")
	metric.HistogramResponseTimeObserve("PostLogin", diffTime.Seconds())
}
