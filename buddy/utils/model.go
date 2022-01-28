package utils

type ForecastResponse struct {
	Forecasts []struct {
		Date        string `json:"date"`
		Weather     string `json:"weather"`
		Temperature struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"temperature"`
		Rain int `json:"rain"`
	} `json:"forecasts"`
}

type WeeklyForecast struct {
	Forecasts []Forecast
}

type Forecast struct {
	Date    int64
	Weather string
	Min     int
	Max     int
	Rain    int
}

type Alarm interface {
	IsRainyTomorrow() bool
	IsThreeDaysRainy() bool
	IsThreeDaysCloudy() bool
	IsFiveDaysSunny() bool
}

func (w *WeeklyForecast) IsRainyTomorrow() bool {
	return w.Forecasts[0].Weather == "rainy"
}

func (w *WeeklyForecast) IsThreeDaysRainy() bool {
	count := 0
	for _, forecast := range w.Forecasts {
		if forecast.Weather == "rainy" {
			count++
		}
	}
	return count >= 3
}

func (w *WeeklyForecast) IsThreeDaysCloudy() bool {
	count := 0
	for _, forecast := range w.Forecasts {
		if forecast.Weather == "cloudy" {
			count++
		}
	}
	return count >= 3
}

func (w *WeeklyForecast) IsFiveDaysSunny() bool {
	count := 0
	for _, forecast := range w.Forecasts {
		if forecast.Weather == "clear" {
			count++
		}
	}
	return count >= 5
}
