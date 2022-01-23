package model

import (
	"sort"
)

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

type VO struct {
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
	Sort()
}

func (v *VO) Sort() {
	sort.Slice(v.Forecasts, func(i, j int) bool {
		return v.Forecasts[i].Date < v.Forecasts[j].Date
	})
}

func (v *VO) IsRainyTomorrow() bool {
	return v.Forecasts[0].Weather == "rainy"
}

func (v *VO) IsThreeDaysRainy() bool {
	count := 0
	for _, forecast := range v.Forecasts {
		if forecast.Weather == "rainy" {
			count++
		}
	}
	return count >= 3
}

func (v *VO) IsThreeDaysCloudy() bool {
	count := 0
	for _, forecast := range v.Forecasts {
		if forecast.Weather == "cloudy" {
			count++
		}
	}
	return count >= 3
}

func (v *VO) IsFiveDaysSunny() bool {
	count := 0
	for _, forecast := range v.Forecasts {
		if forecast.Weather == "clear" {
			count++
		}
	}
	return count >= 5
}
