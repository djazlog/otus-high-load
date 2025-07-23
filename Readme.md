# Монолит работы с пользователем

## Старт проекта
```
make install-deps
make up

go run ./cmd/server/main.go
```

# Импорт данных 
```
go run ./cmd/importer/main.go
```


# Создание миграций
```
goose create -dir migrations create_user_table sql
```
