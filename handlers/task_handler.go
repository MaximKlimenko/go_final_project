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

func AddTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req utils.TaskJSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		// Проверка обязательного поля Title
		if req.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Установка текущей даты, если дата не указана
		if req.Date == "" {
			req.Date = time.Now().Format(utils.TimeFormat)
		}

		// Проверка формата даты
		if _, err := time.Parse(utils.TimeFormat, req.Date); err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
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

		// Добавление задачи в базу данных
		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при добавлении задачи"}`, http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка получения ID задачи"}`, http.StatusInternalServerError)
			return
		}

		// Формирование успешного ответа
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": strconv.FormatInt(id, 10)}
		json.NewEncoder(w).Encode(response)
	}
}

func GetTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task utils.Task
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}
		err := db.Get(&task, `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}
		response := utils.TaskJSON{
			ID:      fmt.Sprint(task.ID), // Преобразуем ID в строку
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(response)
	}
}
