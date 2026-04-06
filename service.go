package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 1. Обновляем структуру Сервиса
type Service struct {
	repo *Repository
	rdb  *redis.Client
}

// 2. Обновляем конструктор
func NewService(repo *Repository, rdb *redis.Client) *Service {
	return &Service{
		repo: repo,
		rdb:  rdb,
	}
}

// 3. Внедряем паттерн Cache-Aside
func (s *Service) GetLink(shortCode string) (string, error) {
	ctx := context.Background()

	// Шаг 1: Ищем в Redis
	val, err := s.rdb.Get(ctx, shortCode).Result()
	if err == nil {
		fmt.Println("⚡ Отдали из кэша (Redis)!")
		return val, nil
	}

	// Шаг 2: Если в кэше нет, идем в Postgres
	fmt.Println("🐢 Идем в базу (Postgres)...")
	longURL, err := s.repo.getLink(shortCode)
	if err != nil {
		return "", err // Ошибка БД или запись не найдена
	}

	// Шаг 3: Нашли в БД? Асинхронно сохраняем в Redis на 24 часа
	// Чтобы не заставлять пользователя ждать сохранения в кэш, можно запустить это в горутине,
	// но пока сделаем синхронно для простоты понимания
	s.rdb.Set(ctx, shortCode, longURL, 24*time.Hour)

	return longURL, nil
}

// ... методы generateShortURL и createShortLink остаются без изменений ...
func generateShortURL(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *Service) createShortLink(originalURL string) (string, error) {
	for {
		shortCode := generateShortURL(6)
		err := s.repo.saveLink(shortCode, originalURL)

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			return "", err
		}
		return shortCode, nil
	}
}