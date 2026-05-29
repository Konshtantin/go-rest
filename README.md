# go-rest

Pet-проект для практики REST API. Postgres c pgx (без ORM намеренно)? Chi-роутер чтобы не уходить в абстракции больших фреймворков типо Gin или Fiber, cache-aside на Redis. Всё поднимается с docker compose.




## Запуск

```
docker compose up --build -d
```

API будет доступен на http://localhost:8080 по умолчанию


## Эндпоинты

### users

- `GET /users` — список пользователей
- `POST /users` — создать нового пользователя
- `GET /users/{id}` — получить по ID (с Redis кэшем)
- `DELETE /users/{id}` — удалить

### posts

- `GET /posts` — список постов
- `POST /posts` — создать пост
- `GET /posts/{id}` — получить по ID
- `DELETE /posts/{id}` — удалить


## Запуск тестов

```bash
go test ./...
```