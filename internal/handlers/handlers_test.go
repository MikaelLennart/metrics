package handlers

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/MikaelLennart/metrics.git/internal/store"
)

func TestGetAllMetrics(t *testing.T) {
	type args struct {
		s *store.MemStorage
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllMetrics(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
