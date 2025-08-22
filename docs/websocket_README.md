# WebSocket Server Documentation

## Архитектура

Приложение теперь использует **раздельную архитектуру** с двумя независимыми серверами:

### 🚀 HTTP Server (REST API)
- **Порт**: 8089 (по умолчанию)
- **Переменные окружения**: `HTTP_HOST`, `HTTP_PORT`
- **Назначение**: REST API для создания постов, управления пользователями и друзьями

### 🔌 WebSocket Server (Real-time Updates)
- **Порт**: 8090 (по умолчанию)  
- **Переменные окружения**: `WEBSOCKET_HOST`, `WEBSOCKET_PORT`
- **Назначение**: WebSocket соединения для real-time обновлений ленты

## Преимущества разделения

1. **Масштабируемость**: Каждый сервер можно масштабировать независимо
2. **Изоляция**: Проблемы в одном сервере не влияют на другой
3. **Производительность**: Оптимизация под конкретные задачи
4. **Безопасность**: Разные порты для разных типов трафика

## Конфигурация

### Переменные окружения

```bash
# HTTP Server
HTTP_HOST=localhost
HTTP_PORT=8089

# WebSocket Server  
WEBSOCKET_HOST=localhost
WEBSOCKET_PORT=8090
```

### Запуск

```bash
# Запуск с переменными окружения
HTTP_HOST=localhost HTTP_PORT=8089 WEBSOCKET_HOST=localhost WEBSOCKET_PORT=8090 go run cmd/server/main.go

# Или через Makefile
make run
```

## WebSocket API

### Подключение

```javascript
// Подключение к WebSocket серверу
const wsUrl = 'ws://localhost:8090/post/feed/posted?token=' + encodeURIComponent(jwtToken);
const ws = new WebSocket(wsUrl);
```

### Аутентификация

WebSocket соединения требуют JWT токен, который можно передать:

1. **URL параметр**: `?token=your-jwt-token`
2. **HTTP заголовок**: `Authorization: Bearer your-jwt-token`

### Сообщения

#### Входящие сообщения (от сервера)

```json
{
  "type": "post",
  "payload": {
    "postId": "uuid",
    "postText": "Текст поста",
    "authorUserId": "uuid"
  }
}
```

## Тестирование

### 1. HTML тест клиент

Откройте `docs/websocket_test.html` в браузере и настройте:
- **Server URL**: `ws://localhost:8090` (WebSocket сервер)
- **JWT Token**: Получите через REST API

### 2. Go тест клиент

```bash
# Сборка клиента
make build-websocket-client

# Запуск теста
make websocket-test
```

### 3. Создание поста и проверка обновлений

1. Создайте пост через REST API: `POST http://localhost:8089/posts`
2. Подключитесь к WebSocket: `ws://localhost:8090/post/feed/posted`
3. Получите real-time уведомление о новом посте

## Мониторинг

### Логи

```bash
# HTTP сервер
2025/08/19 08:31:15 HTTP server starting on localhost:8089

# WebSocket сервер  
2025/08/19 08:31:15 WebSocket server starting on localhost:8090
2025/08/19 08:31:15 WebSocket connection registered: user-id
```

### Метрики

- **HTTP сервер**: Prometheus метрики на порту 2112
- **WebSocket сервер**: Встроенные метрики подключений

## Troubleshooting

### Проблема: Не удается подключиться к WebSocket

**Решение**: Проверьте:
1. WebSocket сервер запущен на правильном порту
2. JWT токен валиден
3. Нет блокировки файрвола

### Проблема: Нет real-time обновлений

**Решение**: Проверьте:
1. WebSocket соединение активно
2. RabbitMQ работает
3. Feed worker запущен

### Проблема: Ошибки аутентификации

**Решение**: Проверьте:
1. JWT токен не истек
2. Токен передан правильно (URL параметр или заголовок)
3. Секретный ключ JWT настроен
