package feed

import (
	"context"
)

// Worker интерфейс для воркера материализации ленты
type Worker interface {
	StartWorker(ctx context.Context) error
	StopWorker(ctx context.Context) error
}
