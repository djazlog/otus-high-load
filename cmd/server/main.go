package main

import (
	"context"
	"log"
	"otus-project/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	go func() {
		err = a.RunPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
