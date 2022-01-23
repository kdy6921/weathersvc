package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestAPIServer(apiKey string, data string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("api_key") != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.URL.Path != "/forecast/daily" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		queryLat := r.URL.Query().Get("lat")
		queryLon := r.URL.Query().Get("lon")
		if queryLat == "" || queryLon == "" {
			http.Error(w, "invalid request parameters", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	}))
}

func TestAPIClientDailyForecast(t *testing.T) {
	t.Run("fetch the daily forecast if the request is valid", func(t *testing.T) {
		testAPIKey := "test-api-key"
		server := createTestAPIServer("test-api-key", `{"timestamp": 0, "code": 1, "min_temp": 10.0, "max_temp": 11.0, "rain": 50}`)
		client := NewAPIClient(server.URL, testAPIKey)
		want := APIClientResponse{
			Timestamp: 0,
			Code:      1,
			MinTemp:   10.0,
			MaxTemp:   11.0,
			Rain:      50,
		}
		got, err := client.GetDailyForecast(0.0, 1.1, 1)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if want != got {
			t.Errorf("Expected response %v, got %v", want, got)
		}
	})

	t.Run("fetch the daily forecast if the parameter is invalid", func(t *testing.T) {
		testAPIKey := "test-api-key"
		server := createTestAPIServer("test-api-key", `{"timestamp": 0, "code": 1, "min_temp": 10.0, "max_temp": 11.0, "rain": 50}`)
		client := NewAPIClient(server.URL, testAPIKey)
		want := APIClientResponse{}
		got, err := client.GetDailyForecast(90.1, 1.1, 1)
		if err == nil {
			t.Errorf("Expected error")
		}

		if want != got {
			t.Errorf("Expected response %v, got %v", want, got)
		}
	})

	t.Run("fetch the daily forecast if the api_key is invalid", func(t *testing.T) {
		testAPIKey := "invalid-api-key"
		server := createTestAPIServer("test-api-key", `{"timestamp": 0, "code": 1, "min_temp": 10.0, "max_temp": 11.0, "rain": 50}`)
		client := NewAPIClient(server.URL, testAPIKey)
		want := APIClientResponse{}
		got, err := client.GetDailyForecast(0.1, 1.1, 1)
		if err == nil {
			t.Errorf("Expected error")
		}

		if want != got {
			t.Errorf("Expected response %v, got %v", want, got)
		}
	})
}
