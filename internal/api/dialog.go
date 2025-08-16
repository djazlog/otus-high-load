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

// GetDialogUserIdList - обработчик GET запроса на /dialog/{user_id}/list
func (i *Implementation) GetDialogUserIdList(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	metric.IncRequestCounter()
	timeStart := time.Now()

	// Получаем ID пользователя из контекста (аутентифицированный пользователь)
	ctx := r.Context()

	fromUserId, err := utils.GetUserFromToken(r)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusUnauthorized), "GetDialogUserIdList")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем список сообщений диалога
	messages, err := i.dialogService.GetDialogList(ctx, *fromUserId, string(userId))
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetDialogUserIdList")
		metric.HistogramResponseTimeObserve("GetDialogUserIdListError", diffTime.Seconds())
		http.Error(w, "Failed to get dialog messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Конвертируем и отправляем ответ
	response := converter.ToDialogMessagesFromService(messages)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetDialogUserIdList")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetDialogUserIdList")
	metric.HistogramResponseTimeObserve("GetDialogUserIdList", diffTime.Seconds())
}

// PostDialogUserIdSend - обработчик POST запроса на /dialog/{user_id}/send
func (i *Implementation) PostDialogUserIdSend(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	metric.IncRequestCounter()
	timeStart := time.Now()

	// Получаем ID пользователя из контекста (аутентифицированный пользователь)
	ctx := r.Context()
	fromUserId, err := utils.GetUserFromToken(r)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusUnauthorized), "GetDialogUserIdList")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим тело запроса
	var requestBody *api.PostDialogUserIdSendJSONBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PostDialogUserIdSend")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if requestBody == nil || requestBody.Text == "" {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PostDialogUserIdSend")
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	// Отправляем сообщение
	err = i.dialogService.SendMessage(ctx, *fromUserId, string(userId), string(requestBody.Text))
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PostDialogUserIdSend")
		metric.HistogramResponseTimeObserve("PostDialogUserIdSendError", diffTime.Seconds())
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PostDialogUserIdSend")
	metric.HistogramResponseTimeObserve("PostDialogUserIdSend", diffTime.Seconds())
}
