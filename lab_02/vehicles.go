package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Auto struct {
	Brand string
	Model string
	Year  int
}

func readCSVFile(filename string) ([]Auto, error) {
	var Cars []Auto

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'

	// Odczytaj pierwszy wiersz jako nagłówek
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}

	// Znajduje indeks kolumny "Year"
	var yearIndex int = -1
	for i, col := range header {
		if strings.TrimSpace(col) == "Year" {
			yearIndex = i
			break
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		year, err := strconv.Atoi(strings.TrimSpace(record[yearIndex]))
		if err != nil {
			return nil, fmt.Errorf("invalid year format in CSV: %v", err)
		}

		car := Auto{
			Brand: strings.TrimSpace(record[0]),
			Model: strings.TrimSpace(record[1]),
			Year:  year,
		}

		if car.Brand != "" && car.Model != "" {
			Cars = append(Cars, car)
		}
	}

	return Cars, nil
}

func sortCarsByBrandInbuilt(cars []Auto) {
	sort.Slice(cars, func(i, j int) bool {
		return strings.ToLower(cars[i].Brand) < strings.ToLower(cars[j].Brand)
	})
}

func partition(cars []Auto, low int, high int) int {
	pivot := strings.ToLower(cars[high].Brand)
	i := low - 1

	for j := low; j < high; j++ {
		if strings.ToLower(cars[j].Brand) < pivot {
			i++
			cars[i], cars[j] = cars[j], cars[i]
		}
	}

	cars[i+1], cars[high] = cars[high], cars[i+1]
	return i + 1
}

func sortCarsByBrandQuicksort(cars []Auto, low int, high int) {
	if low < high {
		pi := partition(cars, low, high)
		sortCarsByBrandQuicksort(cars, low, pi-1)
		sortCarsByBrandQuicksort(cars, pi+1, high)
	}
}

func countCarsByYear(cars []Auto, year int) int {
	sumOfCars := 0
	for _, car := range cars {
		if car.Year == year {
			sumOfCars++
		}
	}
	return sumOfCars
}

func getListOfYears(cars []Auto) []int {
	mapOfYears := map[int]bool{}

	for _, car := range cars {
		mapOfYears[car.Year] = true
	}

	return keys(mapOfYears)
}

func keys(m map[int]bool) []int {
	result := make([]int, 0, len(m)) // Prealokacja pamięci dla wydajności
	for k := range m {
		result = append(result, k)
	}
	return result
}

func main() {
	records, err := readCSVFile("lab_02/all-vehicles-model.csv")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	records2 := make([]Auto, len(records))
	copy(records2, records)

	sortCarsByBrandInbuilt(records)
	sortCarsByBrandQuicksort(records2, 0, len(records2)-1)

	fmt.Println("Sorted using built-in sort:")
	for _, record := range records[:10] {
		fmt.Println(record)
	}

	fmt.Println("\nSorted using quicksort:")
	for _, record := range records2[:10] {
		fmt.Println(record)
	}

	listOfYears := getListOfYears(records)
	sort.Ints(listOfYears)

	for _, year := range listOfYears {
		fmt.Println("Year", year, ":", countCarsByYear(records, year), "cars.")
	}
}
