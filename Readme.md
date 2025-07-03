# Монолит работы с пользователем

## Старт проекта
```
make install-deps
make up

go run ./cmd/server/main.go
```


# Создание миграций
```
goose create -dir migrations create_user_table go
```
