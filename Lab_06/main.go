package main

import (
	"bufio"
	"fmt"
	"lab_06/indicators"
	"lab_06/utils"
	"os"
	"strconv"
	"strings"
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

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. Wskaźnik trendu (SMA)")
		fmt.Println("2. Wskaźnik impetu (RSI)")
		fmt.Println("3. Wskaźnik zmienności (ATR)")
		fmt.Println("0. Wyjście")
		fmt.Print("Wybierz opcję: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Podaj okres SMA: ")
			scanner.Scan()
			period, _ := strconv.Atoi(scanner.Text())
			fmt.Println("SMA:", indicators.SMA(closes, period))
		case "2":
			fmt.Print("Podaj okres RSI: ")
			scanner.Scan()
			period, _ := strconv.Atoi(scanner.Text())
			fmt.Println("RSI:", indicators.RSI(closes, period))
		case "3":
			fmt.Print("Podaj okres ATR: ")
			scanner.Scan()
			period, _ := strconv.Atoi(scanner.Text())
			fmt.Println("ATR:", indicators.ATR(data, period))
		case "0":
			fmt.Println("Koniec programu.")
			return
		default:
			fmt.Println("Nieprawidłowa opcja.")
		}
	}
}
