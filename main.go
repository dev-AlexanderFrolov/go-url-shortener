package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()

	// 1. Инициализируем хранилища (красиво и лаконично)
	db := InitDB()
	defer db.Close()

	rdb := InitRedis()
	defer rdb.Close() // Best practice: закрывать соединения

	// 2. Собираем матрешку зависимостей
	repo := NewRepository(db)
	service := NewService(repo, rdb) // Прокидываем rdb в Сервис
	handler := NewHandler(service)

	// 3. Регистрируем роуты
	http.HandleFunc("POST /api/shorten", handler.ShortenHandler)
	http.HandleFunc("GET /{id}", handler.RedirectHandler)

	fmt.Println("🚀 Сервер стартует на порту 8080...")
	http.ListenAndServe(":8080", nil)
}