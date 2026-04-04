package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	initDB() // Инициализируем базу из repository.go

	http.HandleFunc("POST /api/shorten", ShortenHandler)
	http.HandleFunc("GET /{id}", RedirectHandler)

	fmt.Println("🚀 Сервер стартует на порту 8080...")
	http.ListenAndServe(":8080", nil)
}
