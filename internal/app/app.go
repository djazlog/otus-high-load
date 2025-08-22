package app

import (
	"context"
	"log"
	"net"
	"net/http"
	internalApi "otus-project/internal/api"
	"otus-project/internal/closer"
	"otus-project/internal/config"

	"otus-project/internal/metric"
	"otus-project/internal/model"
	feedHandler "otus-project/internal/service/feed"
	websocketHandler "otus-project/internal/service/websocket"
	"otus-project/pkg/api"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App структура приложения
type App struct {
	serviceProvider  *serviceProvider
	httpServer       *http.Server
	websocketServer  *http.Server
	prometheusServer *http.Server
	websocketHandler *internalApi.WebSocketHandler
	feedWorker       feedHandler.Worker
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
		// Останавливаем WebSocket сервис
		if a.serviceProvider != nil {
			a.serviceProvider.WebSocketService().StopHub(context.Background())
		}
		// Останавливаем воркер материализации ленты
		if a.feedWorker != nil {
			a.feedWorker.StopWorker(context.Background())
		}
		closer.CloseAll()
		closer.Wait()
	}()

	// Запускаем HTTP и WebSocket серверы параллельно
	errChan := make(chan error, 2)

	// Запускаем HTTP сервер
	go func() {
		log.Printf("HTTP server starting on %s", a.serviceProvider.HTTPConfig().Address())
		if err := a.runHTTPServer(); err != nil {
			errChan <- err
		}
	}()

	// Запускаем WebSocket сервер
	go func() {
		log.Printf("WebSocket server starting on %s", a.serviceProvider.WebSocketConfig().Address())
		if err := a.runWebSocketServer(); err != nil {
			errChan <- err
		}
	}()

	// Ждем ошибку от любого из серверов
	return <-errChan
}

// RunPrometheus Run запускает приложение
func (a *App) RunPrometheus() error {
	return a.runPrometheus()
}

// initDeps инициализирует зависимости
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initMetrics,
		a.initServiceProvider,
		a.initWebSocket,
		a.initWebSocketServer,
		a.initHTTPServer,
		a.initPrometheus,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	// Запускаем воркер материализации ленты
	if err := a.feedWorker.StartWorker(ctx); err != nil {
		return err
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

// initMetrics инициализирует Метрики
func (a *App) initMetrics(ctx context.Context) error {
	err := metric.Init(ctx)
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

// initWebSocket инициализирует WebSocket
func (a *App) initWebSocket(ctx context.Context) error {
	// Запускаем WebSocket хаб
	err := a.serviceProvider.WebSocketService().StartHub(ctx)
	if err != nil {
		return err
	}

	// Создаем WebSocket обработчик
	a.websocketHandler = internalApi.NewWebSocketHandler(a.serviceProvider.WebSocketService().GetHub())

	// Создаем воркер материализации ленты
	a.feedWorker = NewFeedWorkerAdapter(a.serviceProvider.FeedService(ctx))

	// Подписываем обработчики на события
	eventBus := a.serviceProvider.EventBus()

	// WebSocket обработчик
	wsEventHandler := websocketHandler.NewEventHandler(a.serviceProvider.WebSocketService())
	eventBus.Subscribe(model.EventTypePostCreated, wsEventHandler.HandlePostCreated)

	// Feed обработчик
	feedEventHandler := feedHandler.NewEventHandler(a.serviceProvider.FeedService(ctx))
	eventBus.Subscribe(model.EventTypePostCreated, feedEventHandler.HandlePostCreated)

	// Потребляем feed events из RabbitMQ и отправляем в конкретные WebSocket-соединения
	if err := a.serviceProvider.QueueClient().ConsumeFeedEvents(ctx, func(ctx context.Context, userID string, ev *model.FeedEvent) error {
		wsPost := &model.WebSocketPost{
			PostID:       ev.PostID,
			PostText:     ev.PostText,
			AuthorUserID: ev.AuthorUserID,
		}
		return a.serviceProvider.WebSocketService().SendPostToUser(ctx, userID, wsPost)
	}); err != nil {
		return err
	}

	return nil
}

// initWebSocketServer инициализирует WebSocket сервер
func (a *App) initWebSocketServer(ctx context.Context) error {
	// Создаем мультиплексор для WebSocket сервера
	mux := http.NewServeMux()

	// Добавляем WebSocket маршрут для канала /post/feed/posted
	mux.HandleFunc("/post/feed/posted", a.websocketHandler.HandleWebSocket)

	a.websocketServer = &http.Server{
		Handler: mux,
		Addr:    a.serviceProvider.WebSocketConfig().Address(),
	}

	return nil
}

// initHTTPServer инициализирует HTTP сервер
func (a *App) initHTTPServer(ctx context.Context) error {
	server := a.serviceProvider.ApiImpl(ctx)

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(server, r)

	// Create middleware for validating tokens.
	mw, err := CreateMiddleware()
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}

	h = mw(h)

	// HTTP сервер только для REST API
	a.httpServer = &http.Server{
		Handler: h,
		Addr:    a.serviceProvider.HTTPConfig().Address(),
	}

	return nil
}

// initPrometheus инициализирует Prometheus сервер
func (a *App) initPrometheus(_ context.Context) error {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
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

// runWebSocketServer запускает WebSocket сервер
func (a *App) runWebSocketServer() error {
	list, err := net.Listen("tcp", a.serviceProvider.WebSocketConfig().Address())
	if err != nil {
		return err
	}

	err = a.websocketServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

// runPrometheus запускает Prometheus сервер
func (a *App) runPrometheus() error {

	log.Printf("Prometheus server is running on %s", "localhost:2112")

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// GetServiceProvider возвращает сервис провайдер
func (a *App) GetServiceProvider() *serviceProvider {
	return a.serviceProvider
}
