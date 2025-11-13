// Package api exposes HTTP handlers for the scheduler application's JSON API.
//
// Handlers are registered by calling Init with a ServeMux. The package reads
// runtime configuration from environment variables.
package api

import (
	"encoding/json"
	"net/http"
	"os"
)

const (
	DefaultMaxTasks = 50
	TaskDateFormat  = "20060102"
)

var todoPassword string

// Init registers the API endpoints on the provided ServeMux.
//
// Endpoints added include:
//   - /api/nextdate
//   - /api/task
//   - /api/task/done
//   - /api/tasks
//   - /api/signin
func Init(mux *http.ServeMux) {

	todoPassword = os.Getenv("TODO_PASSWORD")

	mux.HandleFunc("/api/nextdate", GetNextDateHandler)
	mux.HandleFunc("/api/task", auth(taskHandler))
	mux.HandleFunc("/api/task/done", auth(taskDoneHandler))
	mux.HandleFunc("/api/tasks", auth(tasksHandler))
	mux.HandleFunc("/api/signin", signinHandler)

}

// writeJson serializes 'data' as JSON and writes it to the ResponseWriter with
// the provided HTTP status code. On marshal or write errors it responds with
// HTTP 500 and the underlying error message.
func writeJson(w http.ResponseWriter, data any, code int) {

	responseBody, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	_, err = w.Write(responseBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//endpoint handlers below

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}

}
