package main

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

// NewAPIClient creates a new APIClient
func NewAPIClient(baseURL string, apiKey string) APIClient {
	return APIClient{baseURL: baseURL, apiKey: apiKey}
}

// GetDailyForecast fetches the daily forecast from the third-party weather API service.
func (c APIClient) GetDailyForecast(lat float64, lon float64, dateOffset int) (APIClientResponse, error) {
	return APIClientResponse{}, nil
}
