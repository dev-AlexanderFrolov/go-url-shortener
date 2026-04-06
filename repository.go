package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

// 1. УБРАЛИ var db *sql.DB

// 2. Создаем структуру Репозитория
type Repository struct {
    db *sql.DB
}

// 3. Конструктор для создания Репозитория
func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

// 4. Переделываем InitDB, чтобы она ВОЗВРАЩАЛА подключение, а не писала в глобалку
func InitDB() *sql.DB {
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
    return db // Возвращаем объект базы!
}

// 5. Теперь это МЕТОДЫ структуры Repository (обрати внимание на (r *Repository) перед именем)
func (r *Repository) saveLink(shortCode, originalURL string) error {
    _, err := r.db.Exec("INSERT INTO links (short_code, original_url) VALUES ($1, $2)", shortCode, originalURL)
    return err
}

func (r *Repository) getLink(shortCode string) (string, error) {
    var longURL string
    err := r.db.QueryRow("SELECT original_url FROM links WHERE short_code = $1", shortCode).Scan(&longURL)
    return longURL, err
}