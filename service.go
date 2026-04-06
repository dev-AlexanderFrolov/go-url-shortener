package main

import (
	"math/rand"

	"github.com/lib/pq"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 1. Структура Сервиса, которая "зависит" от Репозитория
type Service struct {
    repo *Repository
}

// 2. Конструктор
func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// GetLink возвращает оригинальную ссылку по хэшу.
// Сейчас он просто дергает БД, но в будущем здесь будет логика проверки кэша (Redis)
func (s *Service) GetLink(shortCode string) (string, error) {
    return s.repo.getLink(shortCode)
}

// Это обычная функция, ей зависимости не нужны, оставляем как есть
func generateShortURL(length int) string {
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

// 3. Делаем метод и используем s.repo вместо прямого вызова
func (s *Service) createShortLink(originalURL string) (string, error) {
    for {
        shortCode := generateShortURL(6)
        err := s.repo.saveLink(shortCode, originalURL) // Вызываем метод ИЗ репозитория

        if err != nil {
            if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
                continue
            }
            return "", err
        }
        return shortCode, nil
    }
}