package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaximKlimenko/go_final_project/utils"
	"github.com/jmoiron/sqlx"
)

func GetTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tasks []utils.Task
		var err error

		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50"
		err = db.Select(&tasks, query)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при выборке задач"}`, http.StatusInternalServerError)
			return
		}

		// Если задач нет, возвращаем пустой слайс
		if tasks == nil {
			tasks = []utils.Task{}
		}

		// Формируем ответ
		response := utils.TasksResponse{
			Tasks: make([]utils.TaskJSON, len(tasks)),
		}

		for i, task := range tasks {
			response.Tasks[i] = utils.TaskJSON{
				ID:      fmt.Sprint(task.ID),
				Date:    task.Date,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error": "Ошибка при формировании ответа"}`, http.StatusInternalServerError)
		}
	}
}
