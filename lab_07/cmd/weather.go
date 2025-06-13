package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"weather-cli/api"
	"weather-cli/models"
	"weather-cli/utils"

	"github.com/olekukonko/tablewriter"
)

var (
	aktualnaCmd = flag.NewFlagSet("aktualna", flag.ExitOnError)
	prognozaCmd = flag.NewFlagSet("prognoza", flag.ExitOnError)
	historiaCmd = flag.NewFlagSet("historia", flag.ExitOnError)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("oczekiwano podkomendy: aktualna, prognoza lub historia")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "aktualna":
		miasto := aktualnaCmd.String("miasto", "", "Nazwa miasta")
		aktualnaCmd.Parse(os.Args[2:])
		if *miasto == "" {
			fmt.Println("Użycie: pogoda aktualna -miasto=<miasto>")
			return
		}
		lat, lon, err := utils.CityToCoords(*miasto)
		if err != nil {
			fmt.Println("Nie znaleziono miasta:", *miasto)
			return
		}
		weather, err := api.GetCurrentWeather(lat, lon)
		if err != nil {
			fmt.Println("Błąd pobierania danych pogodowych:", err)
			return
		}
		PrintWeatherTable([]models.WeatherData{weather})

	case "prognoza":
		miasto := prognozaCmd.String("miasto", "", "Nazwa miasta")
		dni := prognozaCmd.Int("dni", 3, "Liczba dni prognozy")
		prognozaCmd.Parse(os.Args[2:])
		if *miasto == "" {
			fmt.Println("Użycie: pogoda prognoza -miasto=<miasto> [-dni=<dni>]")
			return
		}
		lat, lon, err := utils.CityToCoords(*miasto)
		if err != nil {
			fmt.Println("Nie znaleziono miasta:", *miasto)
			return
		}
		forecast, err := api.GetForecast(lat, lon, *dni)
		if err != nil {
			fmt.Println("Błąd pobierania prognozy:", err)
			return
		}
		PrintWeatherTable(forecast)
		AnalyzeForecast(forecast)
		utils.PlotForecast(forecast)

	case "historia":
		miasto := historiaCmd.String("miasto", "", "Nazwa miasta")
		data := historiaCmd.String("data", "", "Data historyczna")
		historiaCmd.Parse(os.Args[2:])
		if *miasto == "" || *data == "" {
			fmt.Println("Użycie: pogoda historia -miasto=<miasto> -data=<YYYY-MM-DD>")
			return
		}
		lat, lon, err := utils.CityToCoords(*miasto)
		if err != nil {
			fmt.Println("Nie znaleziono miasta:", *miasto)
			return
		}
		history, err := api.GetHistory(lat, lon, *data)
		if err != nil {
			fmt.Println("Błąd pobierania danych historycznych:", err)
			return
		}
		PrintWeatherTable(history)

	default:
		fmt.Println("Nieznana komenda:", os.Args[1])
		os.Exit(1)
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
