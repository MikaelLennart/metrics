package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MikaelLennart/metrics.git/internal/handlers"
	"github.com/MikaelLennart/metrics.git/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	storage := store.NewMemStorage()

	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
	}{
		{
			name:       "Method not allowed",
			method:     http.MethodGet,
			url:        "/update/",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Update Counter",
			method:     http.MethodPost,
			url:        "/update/counter/some/12",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Counter Value",
			method:     http.MethodPost,
			url:        "/update/counter/some/1.2",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Update Gauge",
			method:     http.MethodPost,
			url:        "/update/gauge/some/1.2",
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			h := handlers.UpdateMetrics(storage)

			h.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code, "Status code should match")
		})
	}
}
