package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type APIClientResponse struct {
	Timestamp int64   `json:"timestamp"`
	Code      int     `json:"code"`
	MinTemp   float64 `json:"min_temp"`
	MaxTemp   float64 `json:"max_temp"`
	Rain      int     `json:"rain"`
}

// DailyForecastProvider is an interface that provides daily weather forecast
type DailyForecastProvider interface {
	GetDailyForecast(lat float64, lon float64, dateOffset int) (APIClientResponse, error)
}

// APIClient is API Client for third-party weather API service.
type APIClient struct {
	baseURL string
	apiKey  string
}

var apiPath = "forecast/daily"

// NewAPIClient creates a new APIClient
func NewAPIClient(baseURL string, apiKey string) APIClient {
	return APIClient{baseURL: baseURL, apiKey: apiKey}
}

// GetDailyForecast fetches the daily forecast from the third-party weather API service.
func (c APIClient) GetDailyForecast(lat float64, lon float64, dateOffset int) (APIClientResponse, error) {
	if lat < -90 || lat >= 90 {
		return APIClientResponse{}, fmt.Errorf("lat(%f) is invalid", lat)
	}

	if lon < -180 || lon >= 180 {
		return APIClientResponse{}, fmt.Errorf("lon(%f) is invalid", lon)
	}

	if dateOffset < 0 || dateOffset > 10 {
		return APIClientResponse{}, fmt.Errorf("dateOffset(%d) is invalid", dateOffset)
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return APIClientResponse{}, err
	}
	u.Path = apiPath
	q := u.Query()
	q.Set("lat", fmt.Sprint(lat))
	q.Set("lon", fmt.Sprint(lon))
	q.Set("date_offset", fmt.Sprint(dateOffset))
	q.Set("api_key", c.apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return APIClientResponse{}, err
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIClientResponse{}, err
	}

	var acr APIClientResponse
	err = json.Unmarshal(res, &acr)
	if err != nil {
		return APIClientResponse{}, err
	}

	return acr, nil
}
