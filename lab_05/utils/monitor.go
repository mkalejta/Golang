package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type StopTime struct {
	StopName      string `json:"stopShortName"`
	ArrivalTime   string `json:"arrivalTime"`
	DepartureTime string `json:"departureTime"`
}

// Monitoruje trasę między przystankami
func MonitorRoute(routeID string, stops map[int]string) {
	for {
		date := time.Now().Format("2006-01-02")
		url := fmt.Sprintf("https://ckan2.multimediagdansk.pl/stopTimes?date=%s&routeId=%s", date, routeID)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Błąd pobierania danych:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Błąd HTTP: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
			resp.Body.Close()
			time.Sleep(30 * time.Second)
			continue
		}

		var data struct {
			StopTimes []struct {
				StopId      int    `json:"stopId"`
				ArrivalTime string `json:"arrivalTime"`
			} `json:"stopTimes"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Błąd dekodowania JSON:", err)
			resp.Body.Close()
			time.Sleep(30 * time.Second)
			continue
		}
		resp.Body.Close()

		if len(data.StopTimes) == 0 {
			fmt.Println("Brak danych dla podanej trasy.")
			time.Sleep(30 * time.Second)
			continue
		}

		fmt.Printf("\n[%s] Monitoring trasy %s:\n", time.Now().Format("15:04:05"), routeID)
		for i := 0; i < len(data.StopTimes)-1; i++ {
			from := data.StopTimes[i]
			to := data.StopTimes[i+1]

			fromName := stops[from.StopId]
			toName := stops[to.StopId]

			fromArrivalTime := formatTime(from.ArrivalTime)
			toArrivalTime := formatTime(to.ArrivalTime)

			fmt.Printf("Z %s do %s — przyjazd: %s ➡ przyjazd: %s\n",
				fromName, toName, fromArrivalTime, toArrivalTime)
		}

		fmt.Println("--------------------------------------------------")
		fmt.Println("Monitorowanie zakończone. Czekam na kolejne dane...")

		time.Sleep(30 * time.Second)
	}
}

// Uruchamia monitorowanie równolegle dla wielu tras
func MonitorRoutesConcurrently(routeIDs []string, stops map[int]string) {
	var wg sync.WaitGroup
	for _, routeID := range routeIDs {
		wg.Add(1)
		go func(routeID string) {
			defer wg.Done()
			MonitorRoute(routeID, stops)
		}(routeID)
	}
	wg.Wait()
}

// Pomocnicza funkcja do formatowania czasu
func formatTime(isoTime string) string {
	parsedTime, err := time.Parse("2006-01-02T15:04:05", isoTime)
	if err != nil {
		return "nieznany czas"
	}
	return parsedTime.Format("15:04")
}
