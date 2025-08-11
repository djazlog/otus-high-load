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

// PutFriendDeleteUserId удаляет друга.
// PUT /friend/delete/{user_id}
func (i *Implementation) PutFriendDeleteUserId(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	defer func() {
		diffTime := time.Since(timeStart)
		metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PutFriendDeleteUserId")
		metric.HistogramResponseTimeObserve("PutFriendDeleteUserId", diffTime.Seconds())
	}()

	authId, err := utils.GetUserFromToken(r)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	err = i.friendService.DeleteFriend(context.Background(), *authId, userId)
	if err != nil {
		http.Error(w, "Error delete friend", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}
