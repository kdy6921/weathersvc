package main

// WeeklyForecast is a weekly weather condition forecast. Collection of seven consecutive DailyForecasts
type WeeklyForecast struct {
	Forecasts []DailyForecast `json:"forecasts"`
}

// DailyForecast is a daily weather condition forecast.
type DailyForecast struct {
	Date        string      `json:"date"`
	Weather     string      `json:"weather"`
	Temperature Temperature `json:"temperature"`
	Rain        int         `json:"rain"`
}

// Temperature is a container for min and max temperature
type Temperature struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
