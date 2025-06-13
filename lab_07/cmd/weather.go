package main

import (
	"encoding/json"
	"fmt"
	"os"
	"weather-cli/api"
	"weather-cli/models"
	"weather-cli/utils"

	"github.com/olekukonko/tablewriter"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Użycie:")
		fmt.Println("  pogoda aktualna <miasto>")
		fmt.Println("  pogoda prognoza <miasto> [dni]")
		fmt.Println("  pogoda historia <miasto> [data]")
		return
	}

	komenda := os.Args[1]
	miasto := os.Args[2]
	var dni int
	var data string

	if komenda == "prognoza" && len(os.Args) > 3 {
		dni = utils.ParseInt(os.Args[3], 3)
	}
	if komenda == "historia" && len(os.Args) > 3 {
		data = os.Args[3]
	}

	lat, lon, err := utils.CityToCoords(miasto)
	if err != nil {
		fmt.Println("Nie znaleziono miasta:", miasto)
		return
	}

	switch komenda {
	case "aktualna":
		weather, err := api.GetCurrentWeather(lat, lon)
		if err != nil {
			fmt.Println("Błąd pobierania danych pogodowych:", err)
			return
		}
		PrintWeatherTable([]models.WeatherData{weather})
	case "prognoza":
		forecast, err := api.GetForecast(lat, lon, dni)
		if err != nil {
			fmt.Println("Błąd pobierania prognozy:", err)
			return
		}
		PrintWeatherTable(forecast)
		AnalyzeForecast(forecast)
		utils.PlotForecast(forecast)
	case "historia":
		history, err := api.GetHistory(lat, lon, data)
		if err != nil {
			fmt.Println("Błąd pobierania danych historycznych:", err)
			return
		}
		PrintWeatherTable(history)
	default:
		fmt.Println("Nieznana komenda:", komenda)
	}
}

func PrintWeatherTable(data []models.WeatherData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Data", "Temp", "Wilgotność", "Wiatr", "Ciśnienie", "Opis"})
	for _, d := range data {
		table.Append([]string{
			d.Date, fmt.Sprintf("%.1f°C", d.Temperature), fmt.Sprintf("%d%%", d.Humidity),
			fmt.Sprintf("%.1f km/h", d.Wind), fmt.Sprintf("%dhPa", d.Pressure), d.Description,
		})
	}
	table.Render()
}

func AnalyzeForecast(data []models.WeatherData) {
	// Wczytaj progi z pliku konfiguracyjnego
	f, _ := os.Open("config/thresholds.json")
	defer f.Close()
	var t struct {
		HighTemp float64 `json:"high_temp"`
		LowTemp  float64 `json:"low_temp"`
	}
	json.NewDecoder(f).Decode(&t)
	for _, d := range data {
		if d.Temperature >= t.HighTemp {
			fmt.Printf("UWAGA: Wysoka temperatura (%.1f°C) dnia %s\n", d.Temperature, d.Date)
		}
		if d.Temperature <= t.LowTemp {
			fmt.Printf("UWAGA: Niska temperatura (%.1f°C) dnia %s\n", d.Temperature, d.Date)
		}
	}
}
