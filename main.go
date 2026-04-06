package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
    _ = godotenv.Load()

    // 1. Инициализируем БД
    db := InitDB() 
    defer db.Close() // Best practice: закрыть базу при остановке программы

    // 2. Собираем матрешку зависимостей (Dependency Injection)
    repo := NewRepository(db)       // Репозиторий получает БД
    service := NewService(repo)     // Сервис получает Репозиторий
    handler := NewHandler(service)  // Хэндлер получает Сервис

    // 3. Регистрируем роуты, используя МЕТОДЫ хэндлера
    http.HandleFunc("POST /api/shorten", handler.ShortenHandler)
    http.HandleFunc("GET /{id}", handler.RedirectHandler)

    fmt.Println("🚀 Сервер стартует на порту 8080...")
    http.ListenAndServe(":8080", nil)
}