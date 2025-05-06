package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var customerNames = []string{"Jan", "Anna", "Tomek", "Maria", "Ola"}
var items = []string{"Laptop", "Telefon", "Tablet", "Monitor", "Klawiatura"}

func generateOrders(orderChan chan<- Order, totalOrders int) {
	for i := 1; i <= totalOrders; i++ {
		order := Order{
			ID:           i,
			CustomerName: customerNames[rand.Intn(len(customerNames))],
			Items:        []string{items[rand.Intn(len(items))]},
			TotalAmount:  float64(rand.Intn(500) + 100),
		}
		orderChan <- order
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	close(orderChan)
}

func processOrder(order Order) ProcessResult {
	start := time.Now()
	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)

	success := rand.Float32() > 0.2
	var err error
	if !success {
		err = fmt.Errorf("błąd przetwarzania zamówienia ID %d", order.ID)
	}
	return ProcessResult{
		OrderID:      order.ID,
		CustomerName: order.CustomerName,
		Success:      success,
		ProcessTime:  time.Since(start),
		Error:        err,
	}
}

func worker(id int, orderChan <-chan Order, resultChan chan<- ProcessResult, wg *sync.WaitGroup, retries int) {
	defer wg.Done()
	for order := range orderChan {
		var result ProcessResult
		for i := 0; i <= retries; i++ {
			result = processOrder(order)
			if result.Success {
				break
			}
			fmt.Printf("Worker %d: Ponowna próba przetworzenia zamówienia ID %d\n", id, order.ID)
		}
		resultChan <- result
	}
}

func collectResults(resultChan <-chan ProcessResult, done chan<- bool) {
	var total, success int
	for result := range resultChan {
		total++
		if result.Success {
			success++
		} else {
			fmt.Printf("Nieudane zamówienie: ID %d, błąd: %v\n", result.OrderID, result.Error)
		}
	}
	fmt.Printf("Statystyka: Sukces %.2f%%, Niepowodzenie %.2f%%\n",
		float64(success)/float64(total)*100,
		float64(total-success)/float64(total)*100,
	)
	done <- true
}

func main() {
	totalOrders := 20
	numWorkers := 5
	retries := 2

	orderChan := make(chan Order, totalOrders)
	resultChan := make(chan ProcessResult, totalOrders)
	var wg sync.WaitGroup
	done := make(chan bool)

	go generateOrders(orderChan, totalOrders)

	// Uruchomienie workerów
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, orderChan, resultChan, &wg, retries)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	go collectResults(resultChan, done)

	<-done
	fmt.Println("Symulacja zakończona.")
}
