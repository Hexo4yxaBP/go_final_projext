package db

import (
	"database/sql"
	"fmt"
)

// Task represents a scheduler record stored in the SQLite database.
// Date is stored in YYYYMMDD format.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask inserts a new task into the database and returns the inserted row id.
func AddTask(task *Task) (int64, error) {
	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// GetTask retrieves a task by its integer id (string form) and returns it.
func GetTask(id string) (*Task, error) {

	task := Task{}
	// определите запрос
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=$1`

	row := db.QueryRow(query, id)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask updates a task identified by task.ID. It returns an error when
// the id does not exist or on other DB errors.
func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id`
	res, err := db.Exec(query, sql.Named("id", task.ID), sql.Named("comment", task.Comment), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// DeleteTask removes the task with the provided id from the database.
func DeleteTask(id string) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `DELETE FROM scheduler WHERE id=:id`
	res, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}
