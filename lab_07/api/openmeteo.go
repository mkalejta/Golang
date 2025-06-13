package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather-cli/models"
)

func GetCurrentWeather(lat, lon float64) (models.WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return models.WeatherData{}, err
	}
	defer resp.Body.Close()
	var result struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
			Windspeed   float64 `json:"windspeed"`
			Weathercode int     `json:"weathercode"`
		} `json:"current_weather"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.WeatherData{}, err
	}
	return models.WeatherData{
		Date:        time.Now().Format("2006-01-02"),
		Temperature: result.CurrentWeather.Temperature,
		Wind:        result.CurrentWeather.Windspeed,
		Description: WeatherCodeToDesc(result.CurrentWeather.Weathercode),
	}, nil
}

func GetForecast(lat, lon float64, days int) ([]models.WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&daily=temperature_2m_max,temperature_2m_min,weathercode,precipitation_sum,windspeed_10m_max&forecast_days=%d&timezone=auto", lat, lon, days)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Daily struct {
			Time             []string  `json:"time"`
			TemperatureMax   []float64 `json:"temperature_2m_max"`
			TemperatureMin   []float64 `json:"temperature_2m_min"`
			Weathercode      []int     `json:"weathercode"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
			Windspeed10mMax  []float64 `json:"windspeed_10m_max"`
		} `json:"daily"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	var forecast []models.WeatherData
	for i := range result.Daily.Time {
		forecast = append(forecast, models.WeatherData{
			Date:        result.Daily.Time[i],
			Temperature: result.Daily.TemperatureMax[i],
			Wind:        result.Daily.Windspeed10mMax[i],
			Description: WeatherCodeToDesc(result.Daily.Weathercode[i]),
		})
	}
	return forecast, nil
}

func GetHistory(lat, lon float64, date string) ([]models.WeatherData, error) {
	url := fmt.Sprintf("https://archive-api.open-meteo.com/v1/archive?latitude=%f&longitude=%f&start_date=%s&end_date=%s&daily=temperature_2m_max,temperature_2m_min,weathercode,precipitation_sum,windspeed_10m_max&timezone=auto", lat, lon, date, date)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Daily struct {
			Time             []string  `json:"time"`
			TemperatureMax   []float64 `json:"temperature_2m_max"`
			TemperatureMin   []float64 `json:"temperature_2m_min"`
			Weathercode      []int     `json:"weathercode"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
			Windspeed10mMax  []float64 `json:"windspeed_10m_max"`
		} `json:"daily"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	var history []models.WeatherData
	for i := range result.Daily.Time {
		history = append(history, models.WeatherData{
			Date:        result.Daily.Time[i],
			Temperature: result.Daily.TemperatureMax[i],
			Wind:        result.Daily.Windspeed10mMax[i],
			Description: WeatherCodeToDesc(result.Daily.Weathercode[i]),
		})
	}
	return history, nil
}

func WeatherCodeToDesc(code int) string {
	switch code {
	case 0:
		return "Czyste niebo"
	case 1, 2, 3:
		return "Częściowe zachmurzenie"
	case 45, 48:
		return "Mgła"
	case 51, 53, 55:
		return "Mżawka"
	case 61, 63, 65:
		return "Deszcz"
	case 71, 73, 75:
		return "Śnieg"
	case 80, 81, 82:
		return "Przelotne opady"
	default:
		return "Inne"
	}
}
