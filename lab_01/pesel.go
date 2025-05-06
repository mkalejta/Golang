package main

import (
	"fmt"
	"math/rand"
	"time"
)

// GeneratePESEL: geneuje numer PESEL
// Parametry:
// - birthDate: time.Time: reprezentacja daty urodzenia
// - płeć: znak "M" lub "K"
// Wyjscie:
// Tablica z cyframi numeru PESEL

func GenerujPESEL(birthDate time.Time, gender string) [11]int {

	// tablica zawierajaca kolejne cyfry numeru PESEL
	var cyfryPESEL [11]int

	// konwersja daty na dane skladowe
	year := birthDate.Year()
	month := int(birthDate.Month())
	day := birthDate.Day()

	// losowy numer
	randomSerial := rand.Intn(900) + 100 // 3 cyfrowy losowy numer z zakresu 100-999

	cyfryPESEL[0] = year / 10 % 10
	cyfryPESEL[1] = year % 10

	switch {
	case 1800 <= year && year <= 1899:
		cyfryPESEL[2] = month/10 + 8
		cyfryPESEL[3] = month % 10
	case 1900 <= year && year <= 1999:
		cyfryPESEL[2] = month / 10
		cyfryPESEL[3] = month % 10
	case 2000 <= year && year <= 2099:
		cyfryPESEL[2] = month/10 + 2
		cyfryPESEL[3] = month % 10
	case 2100 <= year && year <= 2199:
		cyfryPESEL[2] = month/10 + 4
		cyfryPESEL[3] = month % 10
	case 2200 <= year && year <= 2299:
		cyfryPESEL[2] = month/10 + 6
		cyfryPESEL[3] = month % 10
	}

	cyfryPESEL[4] = day / 10
	cyfryPESEL[5] = day % 10
	cyfryPESEL[6] = randomSerial / 100
	cyfryPESEL[7] = randomSerial / 10 % 10
	cyfryPESEL[8] = randomSerial % 10

	if gender == "K" {
		cyfryPESEL[9] = rand.Intn(9) * 2 % 10
	} else {
		cyfryPESEL[9] = rand.Intn(9)*2%10 + 1
	}

	wagi := [10]int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
	suma := 0

	for i := 0; i <= 9; i++ {
		suma += (cyfryPESEL[i] * wagi[i]) % 10
	}

	cyfra_kontrolna := 10 - (suma % 10)
	cyfryPESEL[10] = cyfra_kontrolna

	return cyfryPESEL
}

// WeryfikujPESEL: weryfikuje poprawność numeru PESEL
// Parametry:
// - cyfryPESEL: Tablica z cyframi numeru PESEL
// Wyjscie:
//zmienna bool

func WeryfikujPESEL(cyfryPESEL [11]int) bool {

	var czyPESEL bool

	wagi := [10]int{1, 3, 7, 9, 1, 3, 7, 9, 1, 3}
	suma := 0

	for i := 0; i <= 9; i++ {
		suma += (cyfryPESEL[i] * wagi[i]) % 10
	}

	cyfra_kontrolna := 10 - (suma % 10)

	czyPESEL = cyfra_kontrolna == cyfryPESEL[10]

	return czyPESEL
}

// Przykład użycia
func main() {

	birthDate := time.Date(1980, 2, 26, 0, 0, 0, 0, time.UTC)
	pesel := GenerujPESEL(birthDate, "M")

	fmt.Println("Wygenerowany PESEL:", pesel)

	fmt.Println("Czy numer PESEL jest poprawny:", WeryfikujPESEL(pesel))

}
