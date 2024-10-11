package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MaximKlimenko/go_final_project/database"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Print("No .env file found")
	}
}

func main() {
	//conecting db
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connect to database: %s", err)
	}

	//starting server
	fmt.Println("Запускаем сервер")
	mux := http.NewServeMux()

	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	portStr, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		fmt.Println("port doesn't exists")
	}
	// лог-контроль
	fmt.Println(portStr)

	err = http.ListenAndServe(portStr, mux)
	if err != nil {
		panic(err)
	}
}
