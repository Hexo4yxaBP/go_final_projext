# Scheduler (TODO) — Go web app

Small Go scheduler / TODO app: JSON API, SQLite, static UI, Docker.

## Stack

- Go (HTTP API in `pkg/api`, server in `pkg/server`)
- SQLite via `modernc.org/sqlite` + `sqlx`
- JWT auth (`golang-jwt/jwt`)
- Dockerfile + tests in `tests/`

## Run locally

```bash
cp .env.example .env
# edit TODO_PASSWORD in .env
go build -o ./bin/scheduler .
./bin/scheduler
```

Open http://localhost:7540/

## Tests

```bash
go test ./... -v
```

## Docker

```bash
docker build -t scheduler-todo .
docker run -d --name scheduler-todo -p 7540:7540 -v "$(pwd)/data":/app/data scheduler-todo
```

## Status

Portfolio / course final project. Not a production SaaS.
