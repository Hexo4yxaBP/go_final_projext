# go_final_projext — Scheduler (TODO) web app

A small Go-based scheduler / TODO web application with a simple frontend.

What it is
- A minimal web server that serves static files from `web/` and provides a small JSON API (see `pkg/api`).
- Uses an SQLite database stored in `data/scheduler.db` by default.

Why it exists
- This repository is a final project example that demonstrates a small full-stack Go app: server, API, SQLite persistence, and static frontend.

Run locally (development)

1. Prepare environment variables. Copy or edit the `.env` file in the project root. Example:

```env
TODO_PORT=7540
TODO_DBFILE=./data/scheduler.db
TODO_MAXTASKSINRESPONSE=50
TODO_PASSWORD=Mellon
```

2. Build and run with Go (recommended while developing):

```bash
# build
go build -o ./cmd/simpledimple_todolist ./cmd/main.go

# run
./cmd/simpledimple_todolist
```

3. Open the app in your browser:

http://localhost:7540/

Notes:
- The server sets its working directory to the project root (the program calls `os.Chdir("..")` in `cmd/main.go`) so static files resolve from `./web`.
- If you change `TODO_PORT` in `.env`, update the URL accordingly.

Running tests

Run the unit/integration tests provided in the `tests/` folder with:

```bash
go test ./... -v
```

The test configuration values are in `tests/settings.go`. Typical values used by the tests:

```go
var Port = 7540
var DBFile = "../data/scheduler.db"
var FullNextDate = true
var Search = true
var Token = `<example JWT token>`
```

Adjust those values if you need tests to point to a different port or DB file.

Docker

This repo includes a `Dockerfile` that builds the Go binary in a builder stage and places the executable and `web/` into an `ubuntu:latest` runtime image. 

Build the image:

```bash
docker build -t simpledimple_todolist .
```

Run the container and mount the host `data/` (so the SQLite DB persists on the host):

```bash
# from project root; $(pwd) expands to the absolute path of the current directory
docker run -d --name simple_dimple_todolist1 -p 7540:7540 -v "$(pwd)/data":/app/data simpledimple_todolist
```


