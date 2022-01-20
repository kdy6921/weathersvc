package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	provider DailyForecastProvider
}

func NewServer(provider DailyForecastProvider) *Server {
	return &Server{provider: provider}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/weekly" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	queryLat := r.URL.Query().Get("lat")
	queryLon := r.URL.Query().Get("lon")

	if queryLat == "" || queryLon == "" {
		http.Error(w, "invalid request parameters", http.StatusInternalServerError)
		return
	}

	lat, err := strconv.ParseFloat(queryLat, 64)
	if err != nil {
		http.Error(w, "invalid request parameters", http.StatusInternalServerError)
		return
	}
	if lat < -90 || lat > 90 {
		http.Error(w, "invalid request parameters", http.StatusInternalServerError)
		return
	}

	lon, err := strconv.ParseFloat(queryLon, 64)
	if err != nil {
		http.Error(w, "invalid request parameters", http.StatusInternalServerError)
		return
	}
	if lon < -90 || lon > 90 {
		http.Error(w, "invalid request parameters", http.StatusInternalServerError)
		return
	}

	forecast := WeeklyForecast{}
	payload, err := json.Marshal(forecast)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	apiResponse, err := s.provider.GetDailyForecast(lat, lon, 1)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("API Response: %+v\n", apiResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
