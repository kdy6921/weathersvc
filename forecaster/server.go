package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Server struct {
	provider DailyForecastProvider
}

func NewServer(provider DailyForecastProvider) *Server {
	return &Server{provider: provider}
}

func (s *Server) routine(w http.ResponseWriter, lat float64, lon float64, day int, res chan APIClientResponse) {
	apiResponse, err := s.provider.GetDailyForecast(lat, lon, day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res <- apiResponse
}

func (s *Server) callAPI(lat float64, lon float64, w http.ResponseWriter) error {
	days := []int{1, 2, 3, 4, 5, 6, 7}
	clist := []APIClientResponse{}

	c := make(chan APIClientResponse)

	for _, day := range days {
		go s.routine(w, lat, lon, day, c)
	}

	for range days {
		res := <-c
		clist = append(clist, res)
	}
	sort.Slice(clist, func(i, j int) bool {
		return clist[i].Timestamp < clist[j].Timestamp
	})
	dlist := []DailyForecast{}
	for _, v := range clist {
		dforecast := DailyForecast{
			time.Unix(v.Timestamp, 0).Format("2006-01-02"),
			code2msg(v.Code),
			Temperature{
				v.MinTemp,
				v.MaxTemp,
			},
			v.Rain,
		}
		dlist = append(dlist, dforecast)
	}
	forecast := WeeklyForecast{dlist}
	payload, err := json.Marshal(forecast)
	if err != nil {
		return fmt.Errorf("internal error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
	return nil
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

	c := make(chan error, 1)
	go func() {
		c <- s.callAPI(lat, lon, w)
	}()

	select {
	case res := <-c:
		if res != nil {
			http.Error(w, res.Error(), http.StatusInternalServerError)
			return
		}
	case <-time.After(2 * time.Second):
		http.Error(w, "Timeout", http.StatusInternalServerError)
		return
	}
}
