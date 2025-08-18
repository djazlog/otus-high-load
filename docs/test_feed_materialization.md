# Тестирование отложенной материализации ленты

## Подготовка

1. Запустите все сервисы:
```bash
make up
```

2. Примените миграции:
```bash
make migrate
```

3. Запустите основной сервер:
```bash
make run
```

4. В отдельном терминале запустите воркер материализации ленты:
```bash
make feed-worker
```

## Тестирование

### 1. Создание пользователей и дружбы

```bash
# Регистрация пользователя 1
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "User1",
    "secondName": "Test",
    "age": 25,
    "biography": "Test user 1",
    "city": "Moscow",
    "password": "password123"
  }'

# Регистрация пользователя 2
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "User2",
    "secondName": "Test",
    "age": 30,
    "biography": "Test user 2",
    "city": "SPb",
    "password": "password123"
  }'

# Получение токенов
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user1_id",
    "password": "password123"
  }'

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user2_id",
    "password": "password123"
  }'
```

### 2. Добавление дружбы

```bash
# User1 добавляет User2 в друзья
curl -X POST http://localhost:8080/friend/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user1_token" \
  -d '{
    "user_id": "user2_id"
  }'

# User2 принимает заявку
curl -X POST http://localhost:8080/friend/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user2_token" \
  -d '{
    "user_id": "user1_id"
  }'
```

### 3. Создание поста

```bash
# User1 создает пост
curl -X POST http://localhost:8080/post/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer user1_token" \
  -d '{
    "text": "Hello, this is my first post!"
  }'
```

### 4. Проверка материализованной ленты

```bash
# Получение ленты User2 (должен увидеть пост User1)
curl -X GET "http://localhost:8080/feed?user_id=user2_id&offset=0&limit=10" \
  -H "Authorization: Bearer user2_token"
```

### 5. Мониторинг RabbitMQ

1. Откройте RabbitMQ Management: http://localhost:15672
2. Логин: guest, пароль: guest
3. Перейдите в раздел "Queues"
4. Найдите очередь `feed.materialization`
5. Проверьте количество сообщений в очереди

### 6. Проверка WebSocket уведомлений

1. Откройте `docs/websocket_test.html` в браузере
2. Подключитесь с токеном User2
3. Создайте пост от User1
4. Проверьте получение уведомления в реальном времени

## Проверка логов

### Логи основного сервера
```bash
# Должны быть сообщения о планировании обновления ленты
grep "Scheduled feed updates" logs/server.log
```

### Логи воркера
```bash
# Должны быть сообщения об обработке задач
grep "Processed feed update task" logs/worker.log
```

## Проверка базы данных

```sql
-- Проверка материализованных лент
SELECT * FROM materialized_feeds ORDER BY created_at DESC;

-- Проверка заданий
SELECT * FROM feed_jobs ORDER BY created_at DESC;

-- Проверка дружбы
SELECT * FROM friends WHERE status = 'accepted';
```

## Ожидаемые результаты

1. **Создание поста**: Должно быть быстрое (не блокируется обновлением лент)
2. **Материализованная лента**: User2 должен увидеть пост User1 в своей ленте
3. **WebSocket уведомления**: User2 должен получить уведомление в реальном времени
4. **RabbitMQ очередь**: Должны быть задачи в очереди `feed.materialization`
5. **Задания в БД**: Должны быть записи в таблице `feed_jobs` со статусом "completed"

## Устранение неполадок

### Проблема: Нет сообщений в RabbitMQ
- Проверьте, что RabbitMQ запущен: `docker ps | grep rabbitmq`
- Проверьте подключение к RabbitMQ в логах

### Проблема: Воркер не обрабатывает задачи
- Проверьте, что воркер запущен: `ps aux | grep feed_worker`
- Проверьте логи воркера на ошибки

### Проблема: Лента не обновляется
- Проверьте, что дружба установлена: `SELECT * FROM friends WHERE status = 'accepted'`
- Проверьте, что миграции применены: `make local-migration-status`

### Проблема: WebSocket не работает
- Проверьте, что WebSocket сервер запущен
- Проверьте токен авторизации
- Проверьте консоль браузера на ошибки
