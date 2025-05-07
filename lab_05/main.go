package main

import (
	"bufio"
	"fmt"
	"lab_05/models"
	"lab_05/utils"
	"os"
)

func ConvertStopsToMap(stops []models.Stop) map[int]string {
	stopMap := make(map[int]string)
	for _, stop := range stops {
		stopMap[stop.StopId] = stop.StopName
	}
	return stopMap
}

func selectStop(stops []models.Stop) *models.Stop {
	fmt.Print("Wpisz nazwę przystanku: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	matched := utils.SearchStopsByName(stops, input)
	if len(matched) == 0 {
		return nil
	}

	fmt.Println("Wyniki:")
	for i, stop := range matched {
		fmt.Printf("[%d] %s (%s)\n", i, stop.StopName, stop.StopCode)
	}

	fmt.Print("Wybierz numer przystanku: ")
	scanner.Scan()
	var choice int
	fmt.Sscanf(scanner.Text(), "%d", &choice)

	if choice >= 0 && choice < len(matched) {
		return &matched[choice]
	}

	fmt.Println("Nieprawidłowy wybór.")
	return nil
}

func selectRoutes(departures []models.Departure) []string {
	for i, d := range departures {
		fmt.Printf("[%d] Linia %d (%s) — %s\n", i, d.RouteID, d.HeadSign, d.EstimatedTime)
	}

	fmt.Print("Ile linii chcesz monitorować? (1 lub 2): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var numLines int
	fmt.Sscanf(scanner.Text(), "%d", &numLines)

	if numLines == 1 {
		fmt.Print("Wybierz numer linii do monitorowania: ")
		scanner.Scan()
		var choice int
		fmt.Sscanf(scanner.Text(), "%d", &choice)

		if choice >= 0 && choice < len(departures) {
			return []string{fmt.Sprintf("%d", departures[choice].RouteID)}
		}
		fmt.Println("Nieprawidłowy wybór.")
	} else if numLines == 2 {
		fmt.Print("Wybierz numer pierwszej linii do monitorowania: ")
		scanner.Scan()
		var choice1 int
		fmt.Sscanf(scanner.Text(), "%d", &choice1)

		fmt.Print("Wybierz numer drugiej linii do monitorowania: ")
		scanner.Scan()
		var choice2 int
		fmt.Sscanf(scanner.Text(), "%d", &choice2)

		if choice1 >= 0 && choice1 < len(departures) && choice2 >= 0 && choice2 < len(departures) {
			return []string{
				fmt.Sprintf("%d", departures[choice1].RouteID),
				fmt.Sprintf("%d", departures[choice2].RouteID),
			}
		}
		fmt.Println("Nieprawidłowy wybór.")
	}

	return nil
}

func main() {
	stops, err := utils.LoadStopsFromFile("data/stops.json")
	if err != nil {
		fmt.Println("Błąd wczytywania przystanków:", err)
		return
	}

	stopMap := ConvertStopsToMap(stops)

	for {
		selectedStop := selectStop(stops)
		if selectedStop == nil {
			fmt.Println("Nie znaleziono przystanku.")
			return
		}

		departures, err := utils.FetchDepartures(selectedStop.StopId)
		if err != nil {
			fmt.Println("Błąd pobierania odjazdów:", err)
			return
		}

		selectedRoutes := selectRoutes(departures)
		if len(selectedRoutes) == 1 {
			utils.MonitorRoute(selectedRoutes[0], stopMap)
		} else if len(selectedRoutes) == 2 {
			utils.MonitorRoutesConcurrently(selectedRoutes, stopMap)
		} else {
			fmt.Println("Nieprawidłowa liczba linii. Wybierz 1 lub 2.")
		}
	}
}
