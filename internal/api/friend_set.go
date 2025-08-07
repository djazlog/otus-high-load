package api

import (
	"context"
	"net/http"
	"otus-project/internal/metric"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

// PutFriendSetUserId добавляет нового друга к пользователю.
// PUT /friends/set?user_id=123&friend_id=456
func (i *Implementation) PutFriendSetUserId(w http.ResponseWriter, r *http.Request, userId api.UserId) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	friendId := r.URL.Query().Get("friend_id")
	if friendId == "" {
		http.Error(w, "Friend id is empty", http.StatusNotFound)
		return
	}

	err := i.friendService.AddFriend(context.Background(), userId, friendId)
	diffTime := time.Since(timeStart)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PutFriendSetUserId")
		metric.HistogramResponseTimeObserve("PutFriendSetUserId", diffTime.Seconds())
		http.Error(w, "Error add friend", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PutFriendSetUserId")
	metric.HistogramResponseTimeObserve("PutFriendSetUserId", diffTime.Seconds())
}
