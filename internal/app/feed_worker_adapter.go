package app

import (
	"context"
	"otus-project/internal/service/feed"
)

// FeedWorkerAdapter адаптер для преобразования FeedService в FeedWorker
type FeedWorkerAdapter struct {
	feedService feed.Service
}

// NewFeedWorkerAdapter создает новый адаптер
func NewFeedWorkerAdapter(feedService feed.Service) feed.Worker {
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
