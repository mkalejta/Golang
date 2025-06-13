package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	inFile, err := os.Open("data/data.csv")
	if err != nil {
		fmt.Println("Błąd otwierania pliku wejściowego:", err)
		return
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	outFile, err := os.Create("data/converted.csv")
	if err != nil {
		fmt.Println("Błąd tworzenia pliku wyjściowego:", err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Zmień nagłówek na format sample.csv
	writer.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume"})

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Błąd czytania CSV:", err)
		return
	}

	for i, row := range records {
		if i == 0 {
			continue // pomiń nagłówek
		}
		// Wejściowy plik: Date,Close/Last,Volume,Open,High,Low
		date := convertDate(row[0])
		closeVal := clean(row[1])
		volume := clean(row[2])
		open := clean(row[3])
		high := clean(row[4])
		low := clean(row[5])

		writer.Write([]string{date, open, high, low, closeVal, volume})
	}
	fmt.Println("Konwersja zakończona. Wynik zapisano w data/converted.csv")
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "$", ""), ",", "")
}

func convertDate(s string) string {
	// Wejściowy format: MM/DD/YYYY, sample: YYYY-MM-DD
	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return s
	}
	return fmt.Sprintf("20%s-%s-%s", parts[2][2:], pad(parts[0]), pad(parts[1]))
}

func pad(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}
