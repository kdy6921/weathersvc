package utils_test

import (
	"buddy/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func getForecasts(filename string) utils.VO {
	res, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("invalid filename")
	}

	var fr utils.ForecastResponse
	err = json.Unmarshal(res, &fr)
	if err != nil {
		log.Fatal("invalid json")
	}
	var vo utils.VO
	for _, v := range fr.Forecasts {
		t, err := time.Parse("2006-01-02", v.Date)
		if err != nil {
			log.Fatal("invalid json")
		}
		tmp := utils.Forecast{
			Date:    t.Unix(),
			Weather: v.Weather,
			Min:     v.Temperature.Min,
			Max:     v.Temperature.Max,
			Rain:    v.Rain,
		}
		vo.Forecasts = append(vo.Forecasts, tmp)
	}
	return vo
}

func TestVO(t *testing.T) {
	t.Run("Return rainy tomorrow", func(t *testing.T) {
		vo := getForecasts("./response_tmrainy.json")
		if !vo.IsRainyTomorrow() {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Return rainy three days", func(t *testing.T) {
		vo := getForecasts("./response_3rainy.json")
		if vo.IsRainyTomorrow() {
			t.Errorf("Expected false, got true")
		}

		if !vo.IsThreeDaysRainy() {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Return cloudy three days", func(t *testing.T) {
		vo := getForecasts("./response_3cloudy.json")
		if vo.IsRainyTomorrow() {
			t.Errorf("Expected false, got true")
		}

		if vo.IsThreeDaysRainy() {
			t.Errorf("Expected false, got true")
		}

		if !vo.IsThreeDaysCloudy() {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Return clear five days", func(t *testing.T) {
		vo := getForecasts("./response_5sunny.json")
		if vo.IsRainyTomorrow() {
			t.Errorf("Expected false, got true")
		}

		if vo.IsThreeDaysRainy() {
			t.Errorf("Expected false, got true")
		}

		if vo.IsThreeDaysCloudy() {
			t.Errorf("Expected false, got true")
		}

		if !vo.IsFiveDaysSunny() {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Return default", func(t *testing.T) {
		vo := getForecasts("./response.json")
		if vo.IsRainyTomorrow() {
			t.Errorf("Expected false, got true")
		}

		if vo.IsThreeDaysRainy() {
			t.Errorf("Expected false, got true")
		}

		if vo.IsThreeDaysCloudy() {
			t.Errorf("Expected false, got true")
		}

		if vo.IsFiveDaysSunny() {
			t.Errorf("Expected false, got true")
		}
	})
}
