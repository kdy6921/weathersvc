package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type MockDailyForecastProvider struct{}

func (m MockDailyForecastProvider) GetDailyForecast(lat float64, lon float64, dateOffset int) (APIClientResponse, error) {
	return APIClientResponse{
		Timestamp: time.Now().Unix(),
		Code:      0,
		MinTemp:   10.0,
		MaxTemp:   20.0,
		Rain:      0,
	}, nil
}

func TestWeeklyForecastServer(t *testing.T) {
	t.Run("returns the forecast if the request is valid", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weekly", nil)
		query := url.Values{"lat": {"0.0"}, "lon": {"0.0"}}
		request.URL.RawQuery = query.Encode()
		writer := httptest.NewRecorder()
		server := NewServer(MockDailyForecastProvider{})
		server.ServeHTTP(writer, request)
		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, writer.Code)
		}
		var forecast WeeklyForecast
		err := json.NewDecoder(writer.Body).Decode(&forecast)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("returns 400 error if the request parameter lat is not in valid range", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weekly", nil)
		query := url.Values{"lat": {"100.0"}, "lon": {"0.0"}}
		request.URL.RawQuery = query.Encode()
		writer := httptest.NewRecorder()
		server := NewServer(MockDailyForecastProvider{})
		server.ServeHTTP(writer, request)
		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, writer.Code)
		}
	})

	t.Run("returns 400 error if the request parameter lon is not in valid range", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weekly", nil)
		query := url.Values{"lat": {"0.0"}, "lon": {"-250.0"}}
		request.URL.RawQuery = query.Encode()
		writer := httptest.NewRecorder()
		server := NewServer(MockDailyForecastProvider{})
		server.ServeHTTP(writer, request)
		if writer.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, writer.Code)
		}
	})
}
