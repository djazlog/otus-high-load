package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"otus-project/internal/app"
	"syscall"
)

func main() {
	ctx := context.Background()

	// Создаем приложение
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	// Создаем канал для сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем воркер материализации ленты
	feedService := a.GetServiceProvider().FeedService(ctx)
	if err := feedService.StartWorker(ctx); err != nil {
		log.Fatalf("failed to start feed worker: %s", err.Error())
	}

	log.Println("Feed materialization worker started")

	// Ожидаем сигнала завершения
	<-sigChan
	log.Println("Received shutdown signal, stopping feed worker...")

	// Останавливаем воркер
	if err := feedService.StopWorker(ctx); err != nil {
		log.Printf("Error stopping feed worker: %v", err)
	}

	log.Println("Feed materialization worker stopped")
}
