package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	q := "Bengaluru"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=6094ca9b0e9c4a0a9fd83140240604&q=" + q)

	if err != nil {
		panic(err)
	}

	body, errr := io.ReadAll(res.Body)

	if errr != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour

	fmt.Printf("%s %s: %0.1fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		message := fmt.Sprintf("%s - %0.1fC, %0.1f %s", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)

		if hour.ChanceOfRain < 40 {
			fmt.Println(message)
		} else {
			color.Cyan(message)
		}

	}
}
