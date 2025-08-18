# Монолит работы с пользователем

## Старт проекта
```
make install-deps
make up

go run ./cmd/server/main.go
```

## WebSocket Server

Проект включает WebSocket сервер для получения уведомлений о новых постах в реальном времени.

### Подключение к WebSocket

WebSocket сервер доступен по адресу: `ws://localhost:8080/post/feed/posted`

Для подключения требуется JWT токен авторизации.

### Тестирование WebSocket

1. **HTML клиент**: Откройте `docs/websocket_test.html` в браузере
2. **Go клиент**: `go run cmd/websocket_client/main.go ws://localhost:8080 <jwt_token>`

Подробная документация: [docs/websocket_README.md](docs/websocket_README.md)

## Отложенная материализация ленты

Система использует отложенную материализацию ленты через RabbitMQ для оптимизации производительности.

### Запуск воркера материализации

```bash
# Запуск воркера материализации ленты
make feed-worker
```

### Мониторинг RabbitMQ

- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

Подробная документация: [docs/feed_materialization_README.md](docs/feed_materialization_README.md)

# Импорт данных 
```
go run ./cmd/importer/main.go
или
make import
```


# Создание миграций
```
goose create -dir migrations create_user_table sql
goose create -dir migrations materialized_feeds_ефиду sql
```
