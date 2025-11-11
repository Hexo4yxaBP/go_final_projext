package db

import (
	"fmt"
	"time"
)

// Tasks returns up to tasksCount tasks optionally filtered by searchString.
// If searchString parses as a date (DD.MM.YYYY) it's treated as a date filter,
// otherwise it performs a case-insensitive LIKE search against title and comment.
func Tasks(tasksCount int, searchString string) ([]*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1`
	subQuery := ""

	if len(searchString) > 0 {
		searchDate, err := time.Parse("02.01.2006", searchString)
		if err == nil {
			subQuery += fmt.Sprintf(` and date=%s`, searchDate.Format("20060102"))
		} else {
			subQuery += fmt.Sprintf(` and (LOWER(title) LIKE '%%%s%%' or LOWER(comment) LIKE '%%%s%%')`, searchString, searchString)
		}
		query += subQuery
	}

	query += " ORDER BY date LIMIT $1"
	rows, err := db.Query(query, tasksCount)

	if err != nil {
		return nil, err
	}

	resp := []*Task{}

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		resp = append(resp, &t)
	}

	return resp, nil
}
