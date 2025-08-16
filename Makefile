include .env

LOCAL_BIN:=$(CURDIR)/bin
MIGRATION_DIR = migrations

# Для прода
build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/server/main.go

# Запуск сервиса
up:
	@docker compose up -d --build --scale worker=2

# Запуск сервиса для разработки
run:
	go run ./cmd/server/main.go

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
#gen:
#    oapi-codegen \
#    - generate
#    -package spec ./docs/openapi.json > src/gen/gen.go
generate:
	go generate ./...
