package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Драйвер для чтения файлов миграций с диска
	_ "github.com/lib/pq"
)

var db *sql.DB

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

	// --- НАСТРОЙКА МИГРАЦИЙ ---
	// 1. Создаем "драйвер" для работы golang-migrate с нашей базой
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Не удалось создать драйвер миграций: ", err)
	}

	// 2. Указываем, где лежат файлы (file://migrations) и подключаем драйвер БД
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Не удалось инициализировать миграции: ", err)
	}

	// 3. Запускаем накатывание миграций (Up)
	err = m.Up()
	// Ошибка ErrNoChange означает, что новых миграций нет, база актуальна. Это НОРМАЛЬНО.
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Ошибка применения миграций: ", err)
	}

	log.Println("✅ База данных и миграции успешно инициализированы!")
}

// ... функции saveLink и getLink остаются без изменений ...
func saveLink(shortCode, originalURL string) error {
	_, err := db.Exec("INSERT INTO links (short_code, original_url) VALUES ($1, $2)", shortCode, originalURL)
	return err
}

func getLink(shortCode string) (string, error) {
	var longURL string
	err := db.QueryRow("SELECT original_url FROM links WHERE short_code = $1", shortCode).Scan(&longURL)
	return longURL, err
}
