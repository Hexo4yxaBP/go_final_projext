package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Hexo4yxaBP/go_final_projext/pkg/db"
)

// addTaskHandler handles POST /api/task for adding a new task.
// It validates the request body, ensures a title is present, normalizes the
// task date using checkDate and inserts the task into the database.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	//read request body
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(task.Title) == 0 {

		writeJson(w, map[string]string{"error": "title must not be empty"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	taskId, err := db.AddTask(&task)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	task.ID = strconv.FormatInt(taskId, 10)

	writeJson(w, map[string]string{"id": task.ID}, http.StatusCreated)

}

// checkDate ensures a Task has a valid Date field. If the date is empty it
// sets it to today's date. For repeating tasks it computes the next occurrence
// using nextDate when necessary. It returns an error when parsing fails.
func checkDate(task *db.Task) error {

	now, err := time.Parse("20060102", time.Now().Format("20060102"))
	if err != nil {
		return err
	}

	if len(task.Date) == 0 {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	var next string
	if len(task.Repeat) > 0 {
		next, err = nextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if now.After(t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil

}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")

	if len(idString) == 0 {
		writeJson(w, map[string]string{"error": "id parameter is not provided"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(idString)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	//read request body
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(task.Title) == 0 {
		writeJson(w, map[string]string{"error": "title must not be empty"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	response := map[string]string{}

	writeJson(w, response, http.StatusCreated)

}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")

	if len(idString) == 0 {
		writeJson(w, map[string]string{"error": "id parameter is not provided"}, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(idString)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)

}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")

	if len(idString) == 0 {
		writeJson(w, map[string]string{"error": "id parameter is not provided"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(idString)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if len(task.Repeat) == 0 {
		err := db.DeleteTask(idString)

		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		writeJson(w, map[string]string{}, http.StatusOK)
		return
	}

	task.Date, err = nextDate(time.Now(), task.Date, task.Repeat)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	err = db.UpdateTask(task)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)

}
