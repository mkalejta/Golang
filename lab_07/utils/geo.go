package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"weather-cli/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Pobiera współrzędne miasta z API geocode.maps.co
func CityToCoords(city string) (float64, float64, error) {
	url := fmt.Sprintf("https://geocode.maps.co/search?q=%s", strings.ReplaceAll(city, " ", "+"))
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return 0, 0, err
	}
	if len(results) == 0 {
		return 0, 0, fmt.Errorf("nie znaleziono miasta: %s", city)
	}
	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return 0, 0, err
	}
	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return 0, 0, err
	}
	return lat, lon, nil
}

func ParseInt(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

// Prosta wizualizacja wykresu (gonum/plot)
func PlotForecast(data []models.WeatherData) {
	if len(data) == 0 {
		fmt.Println("Brak danych do wykresu.")
		return
	}

	p := plot.New()
	p.Title.Text = "Prognoza temperatury"
	p.X.Label.Text = "Dzień"
	p.Y.Label.Text = "Temperatura (°C)"

	points := make(plotter.XYs, len(data))
	for i, d := range data {
		points[i].X = float64(i)
		points[i].Y = d.Temperature
	}

	err := plotutil.AddLinePoints(p, "Temperatura", points)
	if err != nil {
		fmt.Println("Błąd rysowania wykresu:", err)
		return
	}

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "forecast.png"); err != nil {
		fmt.Println("Błąd zapisu wykresu:", err)
		return
	}
	fmt.Println("Wykres zapisano do pliku forecast.png")
}
