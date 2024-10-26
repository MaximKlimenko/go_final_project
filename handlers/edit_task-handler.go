package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MaximKlimenko/go_final_project/nextdate"
	"github.com/MaximKlimenko/go_final_project/utils"
	"github.com/jmoiron/sqlx"
)

func EditTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.TaskJSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		//Проверка на наличие обязательных полей
		if req.ID == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}
		if req.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Проверка формата даты
		if _, err := time.Parse(utils.TimeFormat, req.Date); err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}
		// Преобразование ID из строки в int64
		id, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Некорректный идентификатор"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли задача
		var existingTask utils.Task
		err = db.Get(&existingTask, `SELECT * FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		now := time.Now().Format(utils.TimeFormat)
		if req.Date < now {
			if req.Repeat == "" {
				req.Date = now
			} else {
				nextDate, err := nextdate.NextDate(time.Now(), req.Date, req.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
					return
				}
				req.Date = nextDate
			}
		}

		//Обновление задачи в бд
		query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
		result, err := db.Exec(query, req.Date, req.Title, req.Comment, req.Repeat, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		// Возврат пустого JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
