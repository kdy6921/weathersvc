package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"buddy/utils"
)

var (
	host    = os.Getenv("API_HOST")
	apiPath = "weekly"
)

func main() {
	pos := utils.GetCurrentGeoPosition()
	u, err := url.Parse(host)
	if err != nil {
		log.Fatal("host is invalid")
	}
	u.Path = apiPath
	q := u.Query()
	q.Set("lat", fmt.Sprint(pos.Latitude))
	q.Set("lon", fmt.Sprint(pos.Longitude))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
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
	var vo utils.WeeklyForecast
	for _, v := range fr.Forecasts {
		t, err := time.Parse("2006-01-02", v.Date)
		if err != nil {
			return
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

	if vo.IsRainyTomorrow() {
		fmt.Print("내일은 비가 내릴 예정입니다!\n")
		return
	}

	if vo.IsThreeDaysRainy() {
		fmt.Print("이번주 내내 비 소식이 있어요.\n")
		return
	}

	if vo.IsThreeDaysCloudy() {
		fmt.Print("날씨가 약간은 칙칙해요\n")
		return
	}

	if vo.IsFiveDaysSunny() {
		fmt.Print("일주일 내내 날씨가 맑아요!\n")
		return
	}

	fmt.Print("맑은 날씨를 즐기세요.\n")
}
