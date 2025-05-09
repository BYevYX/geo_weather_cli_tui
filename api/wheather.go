package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Current struct {
	Time               string  `json:"time"`
	Temperature2m      float64 `json:"temperature_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	WindSpeed10m       float64 `json:"wind_speed_10m"`
}

type CurrentWeather struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Elevation float64 `json:"elevation"`
	Current   Current `json:"current"`
}


func GetCurrentWeather(lat float64, lon float64, timezone string) CurrentWeather {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,apparent_temperature,wind_speed_10m", lat, lon)
	
	if timezone != "" {
		url += fmt.Sprintf("&timezone=%s", timezone)
	}
	resp, err := http.Get(url)
	var weather CurrentWeather

	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("Ошибка чтения ответа:", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("Failed to decode JSON: %v\n", err)
	}

	return weather
}

type Daily struct {
	Time                     []string  `json:"time"`
	Temperature2mMax         []float64 `json:"temperature_2m_max"`
	Temperature2mMin         []float64 `json:"temperature_2m_min"`
	ApparentTemperatureMax   []float64 `json:"apparent_temperature_max"`
	ApparentTemperatureMin   []float64 `json:"apparent_temperature_min"`
	UVIndexMax               []float64 `json:"uv_index_max"`
	WindSpeed10mMax          []float64 `json:"wind_speed_10m_max"`
	WindGusts10mMax          []float64 `json:"wind_gusts_10m_max"`
	WindDirection10mDominant []float64 `json:"wind_direction_10m_dominant"`
	ShortwaveRadiationSum    []float64 `json:"shortwave_radiation_sum"`
	RainSum                 []float64 `json:"rain_sum"`
	SnowfallSum             []float64 `json:"snowfall_sum"`
	ShowersSum              []float64 `json:"showers_sum"`
}

type DailyForecast struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Elevation float64 `json:"elevation"`
	Daily     Daily   `json:"daily"`
}


func GetForecast(lat float64, lon float64, days string, timezone string) DailyForecast {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,uv_index_max,wind_speed_10m_max,wind_gusts_10m_max,wind_direction_10m_dominant,shortwave_radiation_sum,rain_sum,snowfall_sum,showers_sum", lat, lon)
	if timezone != "" {
		url += fmt.Sprintf("&timezone=%s", timezone)
	}
	if days != "" {
		url += fmt.Sprintf("&forecast_days=%s", days)
	}

	resp, err := http.Get(url)
	var weather DailyForecast
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("Ошибка чтения ответа:", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("Failed to decode JSON: %v\n", err)
	}

	return weather
}

