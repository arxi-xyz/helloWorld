package main

import "fmt"

const (
	batchSize = 100
)

type Order string

func main() {
	orders := []Order{"order1", "order2", "order3"}

	batchData := batchProcess(orders)

	fmt.Println("Batch Data:")
	for i, batch := range batchData {
		fmt.Printf("Batch %d: %v\n", i+1, batch)
	}
}

func batchProcess(orders []Order) [][]Order {
	newData := make([]Order, len(orders))
	copy(newData, orders)

	batchData := make([][]Order, 0, (len(newData)+batchSize-1)/batchSize)

	for i := 0; i < len(newData); i += batchSize {
		end := i + batchSize
		if end > len(newData) {
			end = len(newData)
		}

		batch := newData[i:end:end]
		batchData = append(batchData, batch)
	}

	return batchData
}
