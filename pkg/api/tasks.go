package api

import (
	"net/http"

	"github.com/Hexo4yxaBP/go_final_projext/pkg/db"
)

// TasksResp is the JSON response wrapper used by /api/tasks.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler returns a list of tasks, optionally filtered by a search query.
// The maximum number of returned tasks is limited by TODO_MAXTASKSINRESPONSE.
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	searchString := r.URL.Query().Get("search")

	tasks, err := db.Tasks(getMaxTasks(), searchString) // в параметре максимальное количество записей

	if err != nil {
		writeJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}
