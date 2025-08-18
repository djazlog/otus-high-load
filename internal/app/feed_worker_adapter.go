package app

import (
	"context"
	"otus-project/internal/interfaces"
)

// FeedWorkerAdapter адаптер для преобразования FeedService в FeedWorker
type FeedWorkerAdapter struct {
	feedService interfaces.FeedService
}

// NewFeedWorkerAdapter создает новый адаптер
func NewFeedWorkerAdapter(feedService interfaces.FeedService) interfaces.FeedWorker {
	return &FeedWorkerAdapter{
		feedService: feedService,
	}
}

// StartWorker запускает воркер
func (a *FeedWorkerAdapter) StartWorker(ctx context.Context) error {
	return a.feedService.StartWorker(ctx)
}

// StopWorker останавливает воркер
func (a *FeedWorkerAdapter) StopWorker(ctx context.Context) error {
	return a.feedService.StopWorker(ctx)
}
