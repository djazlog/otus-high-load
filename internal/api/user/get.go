package user

import (
	"context"
	"encoding/json"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/pkg/api"
)

// GetUserGetId - получение пользователя по id
func (i *Implementation) GetUserGetId(w http.ResponseWriter, r *http.Request, id api.UserId) {
	w.Header().Set("Content-Type", "application/json")

	userObj, err := i.userService.Get(context.Background(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToUserFromService(userObj)
	// Отправляем объект userObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
