package utils

import (
	"encoding/json"
	"fmt"
	"lab_05/models"
	"net/http"
	"os"
	"strings"
)

// Wczytuje dane przystanków z pliku JSON
func LoadStopsFromFile(filePath string) ([]models.Stop, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("błąd otwierania pliku: %w", err)
	}
	defer file.Close()

	// Struktura odpowiadająca danym w pliku stops.json
	var data map[string]struct {
		LastUpdate string        `json:"lastUpdate"`
		Stops      []models.Stop `json:"stops"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	// Iteracja po mapie i zwrócenie pierwszej niepustej listy przystanków
	for _, entry := range data {
		if len(entry.Stops) > 0 {
			return entry.Stops, nil
		}
	}

	return nil, fmt.Errorf("brak danych o przystankach w pliku JSON")
}

// Wyszukuje przystanki po nazwie
func SearchStopsByName(stops []models.Stop, name string) []models.Stop {
	var matched []models.Stop
	for _, stop := range stops {
		if strings.Contains(strings.ToLower(stop.StopName), strings.ToLower(name)) {
			matched = append(matched, stop)
		}
	}
	return matched
}

// Pobiera dane o odjazdach z API
func FetchDepartures(stopId int) ([]models.Departure, error) {
	url := fmt.Sprintf("http://ckan2.multimediagdansk.pl/departures?stopId=%d", stopId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("błąd pobierania danych: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd HTTP: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Struktura odpowiadająca zagnieżdżonym danym
	var data struct {
		Departures []models.Departure `json:"departures"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("błąd dekodowania JSON: %w", err)
	}

	return data.Departures, nil
}
