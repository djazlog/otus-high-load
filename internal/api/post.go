package api

import (
	"encoding/json"
	"net/http"
	"otus-project/internal/converter"
	"otus-project/internal/metric"
	"otus-project/internal/model"
	"otus-project/internal/utils"
	"otus-project/pkg/api"
	"strconv"
	"time"
)

// PostPostCreate - обработчик POST запроса на /post/create
func (i *Implementation) PostPostCreate(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestCounter()
	timeStart := time.Now()

	// Получаем ID пользователя из контекста (аутентифицированный пользователь)
	ctx := r.Context()
	authorUserId, err := utils.GetUserFromToken(r)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusUnauthorized), "PostPostCreate")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим тело запроса
	var requestBody *api.PostPostCreateJSONBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PostPostCreate")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if requestBody == nil || requestBody.Text == "" {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PostPostCreate")
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	// Создаем пост
	post := &model.Post{
		Text:         &requestBody.Text,
		AuthorUserId: authorUserId,
	}

	postID, err := i.postService.Create(ctx, post)
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PostPostCreate")
		metric.HistogramResponseTimeObserve("PostPostCreateError", diffTime.Seconds())
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем ID созданного поста
	if err := json.NewEncoder(w).Encode(map[string]string{"id": *postID}); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PostPostCreate")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PostPostCreate")
	metric.HistogramResponseTimeObserve("PostPostCreate", diffTime.Seconds())
}

// PutPostDeleteId - обработчик PUT запроса на /post/delete/{id}
func (i *Implementation) PutPostDeleteId(w http.ResponseWriter, r *http.Request, id api.PostId) {
	metric.IncRequestCounter()
	timeStart := time.Now()

	// Получаем ID пользователя из контекста (аутентифицированный пользователь)
	ctx := r.Context()
	userID, err := utils.GetUserFromToken(r)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusUnauthorized), "PutPostDeleteId")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем пост для проверки авторства
	post, err := i.postService.GetByID(ctx, string(id))
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "PutPostDeleteId")
		metric.HistogramResponseTimeObserve("PutPostDeleteIdError", diffTime.Seconds())
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Проверяем, что пользователь является автором поста
	if post.AuthorUserId == nil || *post.AuthorUserId != *userID {
		metric.IncResponseCounter(strconv.Itoa(http.StatusForbidden), "PutPostDeleteId")
		metric.HistogramResponseTimeObserve("PutPostDeleteIdError", diffTime.Seconds())
		http.Error(w, "Forbidden: you can only delete your own posts", http.StatusForbidden)
		return
	}

	// Удаляем пост
	err = i.postService.Delete(ctx, string(id))
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PutPostDeleteId")
		metric.HistogramResponseTimeObserve("PutPostDeleteIdError", diffTime.Seconds())
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PutPostDeleteId")
	metric.HistogramResponseTimeObserve("PutPostDeleteId", diffTime.Seconds())
}

// GetPostGetId - обработчик GET запроса на /post/get/{id}
func (i *Implementation) GetPostGetId(w http.ResponseWriter, r *http.Request, id api.PostId) {
	metric.IncRequestCounter()
	timeStart := time.Now()

	// Получаем пост по ID
	post, err := i.postService.GetByID(r.Context(), string(id))
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "GetPostGetId")
		metric.HistogramResponseTimeObserve("GetPostGetIdError", diffTime.Seconds())
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Конвертируем и отправляем ответ
	response := converter.ToPostFromService(post)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "GetPostGetId")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetPostGetId")
	metric.HistogramResponseTimeObserve("GetPostGetId", diffTime.Seconds())
}

// PutPostUpdate - обработчик PUT запроса на /post/update
func (i *Implementation) PutPostUpdate(w http.ResponseWriter, r *http.Request) {
	metric.IncRequestCounter()
	timeStart := time.Now()

	// Получаем ID пользователя из контекста (аутентифицированный пользователь)
	ctx := r.Context()
	userID, err := utils.GetUserFromToken(r)
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusUnauthorized), "PutPostUpdate")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим тело запроса
	var requestBody *api.PutPostUpdateJSONBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PutPostUpdate")
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if requestBody == nil || requestBody.Id == "" || requestBody.Text == "" {
		metric.IncResponseCounter(strconv.Itoa(http.StatusBadRequest), "PutPostUpdate")
		http.Error(w, "Id and text are required", http.StatusBadRequest)
		return
	}

	// Получаем пост для проверки авторства
	post, err := i.postService.GetByID(ctx, string(requestBody.Id))
	diffTime := time.Since(timeStart)

	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusNotFound), "PutPostUpdate")
		metric.HistogramResponseTimeObserve("PutPostUpdateError", diffTime.Seconds())
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Проверяем, что пользователь является автором поста
	if post.AuthorUserId == nil || *post.AuthorUserId != *userID {
		metric.IncResponseCounter(strconv.Itoa(http.StatusForbidden), "PutPostUpdate")
		metric.HistogramResponseTimeObserve("PutPostUpdateError", diffTime.Seconds())
		http.Error(w, "Forbidden: you can only update your own posts", http.StatusForbidden)
		return
	}

	// Обновляем пост
	err = i.postService.Update(ctx, string(requestBody.Id), string(requestBody.Text))
	if err != nil {
		metric.IncResponseCounter(strconv.Itoa(http.StatusInternalServerError), "PutPostUpdate")
		metric.HistogramResponseTimeObserve("PutPostUpdateError", diffTime.Seconds())
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "PutPostUpdate")
	metric.HistogramResponseTimeObserve("PutPostUpdate", diffTime.Seconds())
}

// GetPostFeed GET /post/feed
func (i *Implementation) GetPostFeed(w http.ResponseWriter, r *http.Request, params api.GetPostFeedParams) {
	metric.IncRequestCounter()
	w.Header().Set("Content-Type", "application/json")
	timeStart := time.Now()

	defer func() {
		diffTime := time.Since(timeStart)
		metric.IncResponseCounter(strconv.Itoa(http.StatusOK), "GetPostFeed")
		metric.HistogramResponseTimeObserve("GetPostFeed", diffTime.Seconds())
	}()

	userId, err := utils.GetUserFromToken(r)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Convert *float32 Offset and Limit to int, defaulting to 0 if nil
	var offset, limit int
	if params.Offset != nil {
		offset = int(*params.Offset)
	}
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	//postsObj, err := i.postService.Feed(r.Context(), *userId, params.Offset, params.Limit)
	postsObj, err := i.feedService.GetMaterializedFeed(r.Context(), *userId, offset, limit)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Преобразуем postsObj из []*MaterializedFeed в []*Post перед передачей в ToPostsFromService
	posts := make([]*model.Post, 0, len(postsObj))
	for _, mf := range postsObj {
		post := converter.MaterializedFeedToPost(mf)
		posts = append(posts, post)
	}

	w.WriteHeader(http.StatusOK)
	response := converter.ToPostsFromService(posts)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	

}
