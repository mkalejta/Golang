package main

import (
	"fmt"
	"lab_06/indicators"
	"lab_06/utils"
	"os"
)

func main() {
	data, err := utils.LoadCSV("data/converted.csv")
	if err != nil {
		fmt.Println("Błąd wczytywania danych:", err)
		os.Exit(1)
	}

	var closes []float64
	for _, d := range data {
		closes = append(closes, d.Close)
	}

	fmt.Println("Wskaźnik trendu (SMA 5):", indicators.SMA(closes, 5))
	fmt.Println("Wskaźnik impetu (RSI 14):", indicators.RSI(closes, 14))
	fmt.Println("Wskaźnik zmienności (ATR 14):", indicators.ATR(data, 14))
}
