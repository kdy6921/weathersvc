package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
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
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(queryLat, 64)
	if err != nil {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}
	if lat < -90 || lat >= 90 {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(queryLon, 64)
	if err != nil {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}
	if lon < -180 || lon >= 180 {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	dlist := []DailyForecast{}

	for i := 1; i <= 7; i++ {
		apiResponse, err := s.provider.GetDailyForecast(lat, lon, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("API Response: %+v\n", apiResponse)

		dforecast := DailyForecast{
			time.Unix(apiResponse.Timestamp, 0).Format("2006-01-02"),
			code2msg(apiResponse.Code),
			Temperature{
				apiResponse.MinTemp,
				apiResponse.MaxTemp,
			},
			apiResponse.Rain,
		}
		dlist = append(dlist, dforecast)
	}

	forecast := WeeklyForecast{dlist}
	payload, err := json.Marshal(forecast)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
