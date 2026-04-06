package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type ShortenRequest struct {
    URL string `json:"url"`
}

type ShortenResponse struct {
    ShortURL string `json:"short_url"`
}

// 1. Структура Хэндлера
type Handler struct {
    service *Service
}

// 2. Конструктор
func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// 3. Делаем методы из функций
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
    var req ShortenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
        http.Error(w, "Неверный запрос", http.StatusBadRequest)
        return
    }

    // Вызываем сервис через h.service
    shortURL, err := h.service.createShortLink(req.URL)
    if err != nil {
        http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
        return
    }

    res := ShortenResponse{ShortURL: shortURL}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(res)
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
    shortCode := r.PathValue("id")

    // Чисто, красиво и по архитектурным канонам: Хэндлер -> Сервис
    longURL, err := h.service.GetLink(shortCode) 

    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Не найдено", http.StatusNotFound)
            return
        }
        http.Error(w, "Ошибка БД", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, longURL, http.StatusFound)
}