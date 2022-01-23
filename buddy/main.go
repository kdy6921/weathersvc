package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kdy6921/weathersvc/buddy/utils"
)

func main() {
	pos := utils.GetCurrentGeoPosition()
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/weekly?lat=%f&lon=%f", pos.Latitude, pos.Longitude))
	if err != nil {
		log.Fatal("API is broken")
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Response is broken")
	}

	var fr utils.ForecastResponse
	err = json.Unmarshal(res, &fr)
	if err != nil {
		return
	}
	var vo utils.VO
	for _, v := range fr.Forecasts {
		t, err := time.Parse("2006-01-02", v.Date)
		if err != nil {
			return
		}
		tmp := utils.Forecast{
			t.Unix(),
			v.Weather,
			v.Temperature.Min,
			v.Temperature.Max,
			v.Rain,
		}
		vo.Forecasts = append(vo.Forecasts, tmp)
	}
	vo.Sort()

	if vo.IsRainyTomorrow() {
		fmt.Print("내일은 비가 내릴 예정입니다!")
		return
	}

	if vo.IsThreeDaysRainy() {
		fmt.Print("이번주 내내 비 소식이 있어요.")
		return
	}

	if vo.IsThreeDaysRainy() {
		fmt.Print("날씨가 약간은 칙칙해요")
		return
	}

	if vo.IsThreeDaysRainy() {
		fmt.Print("일주일 내내 날씨가 맑아요!")
		return
	}

	fmt.Print("맑은 날씨를 즐기세요.")
}
