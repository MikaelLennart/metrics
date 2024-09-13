package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MikaelLennart/metrics.git/internal/handlers"
	"github.com/MikaelLennart/metrics.git/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	// Создаем новый экземпляр хранения
	storage := store.NewMemStorage()
	r := chi.NewRouter()

	// Настраиваем роутинг
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics(storage))

	tests := []struct {
		name               string
		requestBody        string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Update gauge IsOK",
			requestBody:        "",
			url:                "/update/gauge/testMetric/12",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Update gauge Invalid value",
			requestBody:        "",
			url:                "/update/gauge/testMetric/xxxx",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Update countert IsOK",
			requestBody:        "",
			url:                "/update/counter/somecounter/11",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid metric type",
			requestBody:        "",
			url:                "/update/invalidType/testMetric/42",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Missing value",
			requestBody:        "",
			url:                "/update/gauge/some/",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем HTTP-запрос
			req, err := http.NewRequest("POST", tc.url, strings.NewReader(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Создаем новый HTTP-респондер
			rec := httptest.NewRecorder()

			// Выполняем запрос
			r.ServeHTTP(rec, req)

			// Проверяем статус-код ответа
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
}
