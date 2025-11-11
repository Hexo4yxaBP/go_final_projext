// Command main provides the entrypoint for the scheduler web application.
//
// It performs a small amount of startup wiring:
//  - changes working directory to the repository root so static files under ./web are served
//  - loads environment variables from a local .env file (if present)
//  - initializes the SQLite database via pkg/db
//  - starts the HTTP server returned by pkg/server
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Hexo4yxaBP/go_final_projext/pkg/db"
	"github.com/Hexo4yxaBP/go_final_projext/pkg/server"
	"github.com/joho/godotenv"
)

const (
	DEFAULT_PORT    = "7540"                // Default port to listen on if TODO_PORT is not set.
	DEFAULT_DB_FILE = "./data/scheduler.db" // Default database file path if TODO_DBFILE is not set.
)

// main initializes application resources and starts the HTTP server.
//
// It reads TODO_DBFILE from the environment (falling back to DEFAULT_DB_FILE),
// initializes the DB, then constructs and runs the server. The server's
// ListenAndServe is blocking until the process is stopped.
func main() {
	// Make working directory the parent directory so ./web resolves when starting from cmd/.
	err := os.Chdir("..")

	if err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	//db init
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = DEFAULT_DB_FILE
	}

	err = db.Init(dbFile)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	defer db.Close()

	//server init
	srv, err := server.InitServer(log.Default())

	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(srv.Server.Addr, srv.Server.Handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
