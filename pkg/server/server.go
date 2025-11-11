// Package server wires HTTP handlers and exposes a small server helper used by cmd/main.
//
// The server serves static files from ./web and mounts API endpoints defined in pkg/api.
package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hexo4yxaBP/go_final_projext/pkg/api"
)

const (
	DEFAULT_PORT = "7540" // Default port to listen on if TODO_PORT is not set.
)

type MyServer struct {
	Logger *log.Logger
	Server *http.Server
}

// InitServer creates and configures an HTTP server instance.
//
// The returned *MyServer contains the configured http.Server and a logger. The
// function reads the TODO_PORT environment variable (falling back to DEFAULT_PORT)
// and mounts a file server for ./web and API handlers from pkg/api.
func InitServer(l *log.Logger) (*MyServer, error) {

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	r := http.NewServeMux()

	// Serve static files from the web/ directory. This will serve index.html for the root path
	// and any JS/CSS files under /js/ and /css/.
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/", fs)

	// Serve API endpoints
	api.Init(r)

	srv := &http.Server{
		Addr:         ":" + port, //config.Server.Addr,
		Handler:      r,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,  //config.Server.ReadTimeout,
		WriteTimeout: 10 * time.Second, //config.Server.WriteTimeout,
		IdleTimeout:  15 * time.Second, //config.Server.IdleTimeout,
	}

	mySrv := &MyServer{
		Logger: l,
		Server: srv,
	}

	return mySrv, nil

}
