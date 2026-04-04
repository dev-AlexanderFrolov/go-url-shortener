package main

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

// initDB подключается к базе и создает таблицы
func initDB() {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка БД: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("База недоступна: ", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links (short_code VARCHAR(10) PRIMARY KEY, original_url TEXT NOT NULL)`)
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}
}

// saveLink сохраняет ссылку, возвращает ошибку, если такой код уже есть
func saveLink(shortCode, originalURL string) error {
	_, err := db.Exec("INSERT INTO links (short_code, original_url) VALUES ($1, $2)", shortCode, originalURL)
	return err
}

// getLink ищет оригинальный URL по короткому коду
func getLink(shortCode string) (string, error) {
	var longURL string
	err := db.QueryRow("SELECT original_url FROM links WHERE short_code = $1", shortCode).Scan(&longURL)
	return longURL, err
}
