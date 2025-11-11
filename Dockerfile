# Multi-stage Dockerfile: build with golang, run with ubuntu:latest
FROM golang:1.24.3 AS builder
WORKDIR /src

# Copy source and build a static binary
COPY . .
RUN go mod download


COPY . .
# Build statically where possible (CGO disabled). Adjust if your sqlite driver requires CGO.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags='-s -w' -o /app/cmd/todo ./cmd/main.go


FROM ubuntu:latest
#RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

# Create app dirs. The program calls os.Chdir("..") in cmd/main.go so we place the binary
# in /app/cmd and the web/ directory at /app/web so the working dir becomes /app when it runs.
WORKDIR /app/cmd

# Copy binary and static assets from builder
COPY --from=builder /app/cmd/todo /app/cmd/todo
COPY --from=builder /src/web /app/web

# Expose the default port used by the app
EXPOSE 7540

# Defaults for runtime (can be overridden with --env or env_file)
ENV TODO_PORT=7540 \
	TODO_DBFILE=./data/scheduler.db \
	TODO_PASSWORD=Mellon

CMD ["/app/cmd/todo"]
