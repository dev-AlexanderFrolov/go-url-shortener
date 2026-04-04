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

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	// Вызываем наш Service
	shortURL, err := createShortLink(req.URL)
	if err != nil {
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	res := ShortenResponse{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("id")

	// Вызываем наш Repository
	longURL, err := getLink(shortCode)

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
