package api

import (
	"context"
	"net/http"
	"otus-project/internal/metric"
	"otus-project/internal/utils"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

// PutFriendSetUserId добавляет нового друга к пользователю.
// PUT /friend/set/{user_id}
func (i *Implementation) PutFriendSetUserId(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	defer func() {
		diffTime := time.Since(timeStart)
		metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PutFriendSetUserId")
		metric.HistogramResponseTimeObserve("PutFriendSetUserId", diffTime.Seconds())
	}()

	authId, err := utils.GetUserFromToken(r)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	err = i.friendService.AddFriend(context.Background(), *authId, userId)
	if err != nil {
		http.Error(w, "Error add friend", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
