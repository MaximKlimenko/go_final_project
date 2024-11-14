package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/MaximKlimenko/go_final_project/database"
	"github.com/MaximKlimenko/go_final_project/handlers"
)

func main() {
	//Подключение файла окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	//Подключение к бд
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connect to database: %s", err)
	}

	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = "7540"
	}

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(db))
	http.HandleFunc("/api/task/done", handlers.DoneTaskHandler(db))
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddTaskHandler(db)(w, r)
		case http.MethodGet:
			handlers.GetTaskHandler(db)(w, r)
		case http.MethodPut:
			handlers.EditTaskHandler(db)(w, r)
		case http.MethodDelete:
			handlers.DeleteTaskHandler(db)(w, r)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})
	//Запуск сервера
	serverAddress := fmt.Sprintf("localhost:%s", port)
	log.Println("Listening on " + serverAddress)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
