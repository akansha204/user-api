# User API

REST API in Go for managing users with `name` and `dob`, plus dynamic age calculation on read endpoints.

## Stack

- GoFiber
- PostgreSQL
- SQLC
- Uber Zap
- go-playground/validator

## Project Structure

- `cmd/server/main.go`
- `config/`
- `db/migrations/`
- `db/sqlc/`
- `internal/handler/`
- `internal/repository/`
- `internal/service/`
- `internal/routes/`
- `internal/middleware/`
- `internal/models/`
- `internal/logger/`

## Setup

1. Set environment variables as needed:

```bash
export PORT=:3000
export DB_DRIVER=postgres
export DB_SOURCE='postgres://postgres:postgres@localhost:5432/user_api?sslmode=disable'
```

2. Run database migration(s).
3. Generate SQLC code if you change the queries.
4. Start the server:

```bash
go run ./cmd/server
```

## Docker

```bash
docker compose up --build
```

This starts PostgreSQL and the API container.

## Pagination

`GET /users` supports:

- `page` default `1`
- `limit` default `10`
- max `limit` `100`

Example:

```bash
GET /users?page=2&limit=20
```

## Middleware

- `X-Request-Id` is injected into responses
- request duration is logged with Zap

## API

### Create user

`POST /users`

```json
{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get user

`GET /users/:id`

### Update user

`PUT /users/:id`

### Delete user

`DELETE /users/:id`

### List users

`GET /users`
