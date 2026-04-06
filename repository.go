package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9" // <-- Добавили импорт Redis
)

// Структура и InitDB остаются без изменений...
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func InitDB() *sql.DB {
	// ... твой код InitDB ...
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка БД: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("База недоступна: ", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Не удалось создать драйвер миграций: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Не удалось инициализировать миграции: ", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Ошибка применения миграций: ", err)
	}

	log.Println("✅ База данных и миграции успешно инициализированы!")
	return db
}

// НОВАЯ ФУНКЦИЯ: Инициализация Redis "в одну строчку" для main.go
func InitRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379" // Фолбек, если забыли добавить в .env
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Даем Redis 5 секунд на ответ, иначе падаем
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("❌ Redis недоступен: ", err)
	}

	log.Println("✅ Redis успешно подключен!")
	return rdb
}

func (r *Repository) saveLink(shortCode, originalURL string) error {
	_, err := r.db.Exec("INSERT INTO links (short_code, original_url) VALUES ($1, $2)", shortCode, originalURL)
	return err
}

func (r *Repository) getLink(shortCode string) (string, error) {
	var longURL string
	err := r.db.QueryRow("SELECT original_url FROM links WHERE short_code = $1", shortCode).Scan(&longURL)
	return longURL, err
}