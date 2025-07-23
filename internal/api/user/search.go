package user

import (
	"context"
	"encoding/json"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/pkg/api"
)

func (i *Implementation) GetUserSearch(w http.ResponseWriter, r *http.Request, params api.GetUserSearchParams) {
	w.Header().Set("Content-Type", "application/json")

	filter := converter.ToUserFilterFromApi(&params)

	usersObj, err := i.userService.Search(context.Background(), filter)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToUsersFromService(usersObj)

	// Отправляем объект usersObj в формате JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
