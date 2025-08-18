include .env

LOCAL_BIN:=$(CURDIR)/bin
MIGRATION_DIR = migrations

# Для прода
build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/server/main.go

# Запуск сервиса
up:
	@docker compose up -d --build --scale worker=7

# Запуск сервиса для разработки
run:
	go run ./cmd/server/main.go

# Запуск ребалансировки citus
rebalance:
	go run ./cmd/rebalance/main.go

# Установка зависимостей
install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

# Применение миграций
migrate:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v

# Импорт данных
import:
	go run ./cmd/importer/main.go

# Генерация постов для тестирования
posts:
	go run ./cmd/feel_posts/main.go

# WebSocket тестирование
websocket-test:
	@echo "Для тестирования WebSocket используйте:"
	@echo "1. HTML клиент: открыть docs/websocket_test.html в браузере"
	@echo "2. Go клиент: go run cmd/websocket_client/main.go ws://localhost:8080 <jwt_token>"
	@echo ""
	@echo "Сначала получите JWT токен через REST API /login"

# Сборка WebSocket клиента
build-websocket-client:
	go build -o bin/websocket_client cmd/websocket_client/main.go

# Запуск воркера материализации ленты
feed-worker:
	go run ./cmd/feed_worker/main.go

# Сборка воркера материализации ленты
build-feed-worker:
	go build -o bin/feed_worker cmd/feed_worker/main.go

#gen:
#    oapi-codegen \
#    - generate
#    -package spec ./docs/openapi.json > src/gen/gen.go
generate:
	go generate ./...
