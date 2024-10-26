package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MaximKlimenko/go_final_project/utils"
	"github.com/jmoiron/sqlx"
)

func GetTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []utils.Task
		var err error
		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50"
		rows, err := db.Query(query)
		for rows.Next() {
			var task utils.Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				http.Error(w, `{"error":"Ошибка при чтении результатов"}`, http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, `{"error":"Ошибка при обработке результатов"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
	}
}
