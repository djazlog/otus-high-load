package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"otus-project/internal/closer"
	"otus-project/internal/config"
	"otus-project/pkg/api"
)

// App структура приложения
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

// NewApp создает новый экземпляр приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHTTPServer()
}

// initDeps инициализирует зависимости
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initConfig инициализирует конфигурацию
func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

// initServiceProvider инициализирует сервис провайдер
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// initHTTPServer инициализирует HTTP сервер
func (a *App) initHTTPServer(ctx context.Context) error {
	server := a.serviceProvider.UserImpl(ctx)

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(server, r)

	a.httpServer = &http.Server{
		Handler: h,
		Addr:    a.serviceProvider.HTTPConfig().Address(),
	}

	return nil
}

// runHTTPServer запускает HTTP сервер
func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.HTTPConfig().Address())
	if err != nil {
		return err
	}

	err = a.httpServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
