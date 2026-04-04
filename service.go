package main

import (
	"math/rand"

	"github.com/lib/pq"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateShortURL генерирует строку
func generateShortURL(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// createShortLink объединяет генерацию и сохранение (разрешает коллизии)
func createShortLink(originalURL string) (string, error) {
	for {
		shortCode := generateShortURL(6)
		err := saveLink(shortCode, originalURL)

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue // Коллизия, пробуем еще раз
			}
			return "", err // Реальная ошибка базы
		}
		return shortCode, nil // Успех!
	}
}
