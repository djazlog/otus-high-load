package interfaces

import (
	"context"
)

// FeedWorker интерфейс для воркера материализации ленты
type FeedWorker interface {
	StartWorker(ctx context.Context) error
	StopWorker(ctx context.Context) error
}
