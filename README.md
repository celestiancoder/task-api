# Task API

A small REST API for managing tasks, built with Go, PostgreSQL, and `pgx`.

## Features

- Create, read, update, and delete tasks
- Filter tasks by completion status
- Search tasks by name
- Paginate and sort task results

## Requirements

- Go 1.26+
- PostgreSQL

## Setup

1. Create a PostgreSQL database.
2. Run the migration in `db/migrations/000001_create_tasks_table.up.sql`.
3. Update the PostgreSQL connection string in `cmd/api/main.go`.
4. Start the API:

```bash
go run ./cmd/api
```

The server runs at `http://localhost:8080`.

## Endpoints

| Method   | Endpoint      | Description   |
| -------- | ------------- | ------------- |
| `POST`   | `/tasks`      | Create a task |
| `GET`    | `/tasks`      | List tasks    |
| `GET`    | `/tasks/{id}` | Get one task  |
| `PUT`    | `/tasks/{id}` | Update a task |
| `DELETE` | `/tasks/{id}` | Delete a task |

`GET /tasks` supports `page`, `limit`, `completed`, `search`, `sort`, and `order` query parameters.

Example request body:

```json
{
  "name": "Learn Go"
}
```

## Test

```bash
go test ./...
```
