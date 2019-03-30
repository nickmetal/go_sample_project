package main

import (
	"encoding/json"
	"fmt"
	"math"
)

type Cell struct {
	x      int     // row index
	y      int     // column index
	value  float64 // how much goods me provided //transporting price from producers to consumers
	marked bool    // mark as processed by NorthWest Method
}

type Message struct {
	Prices           [][]float64 `json:"prices"`
	ProducersSources []float64   `json:"producers_sources"`
	ConsumersNeeds   []float64   `json:"consumers_needs"`
}

type PriceMatrix [][]Cell

/*
	TODOs
	- check is input data do closed type problem (all need and sources are equaled)
	- check counts needs and sources with prices col and rows num

*/
func initPriceMatrix(m Message) *PriceMatrix {
	pm := make(PriceMatrix, 0)
	for rowIdx, row := range m.Prices {
		pm = append(pm, []Cell{})
		for columnIdx, cellValue := range row {
			cell := Cell{x: rowIdx, y: columnIdx, value: cellValue}
			pm[rowIdx] = append(pm[rowIdx], cell)
		}
	}
	return &pm
}

func findBasicSolution(pm *PriceMatrix, m *Message) {

	for rowIdx, row := range m.Prices {
		// for columnIdx, cellValue := range row {
		for columnIdx, _ := range row {
			matrix := *(pm)
			inputs := *(m)

			if matrix[rowIdx][columnIdx].marked {
				continue
			}

			matrix[rowIdx][columnIdx].value = math.Min(inputs.ProducersSources[rowIdx], inputs.ConsumersNeeds[columnIdx])

			// if we provide full needs
			if inputs.ConsumersNeeds[columnIdx]-matrix[rowIdx][columnIdx].value == 0 {
				fmt.Printf("found full: %v \n", matrix[rowIdx][columnIdx])
				break
			}
		}
	}
}

func printPM(pm PriceMatrix) {
	fmt.Println("priceMatrix")
	for _, row := range pm {
		fmt.Println(row)
	}
}

func main() {
	msg := []byte(`{
		"consumers_needs": [20, 30, 30, 10], 
		"producers_sources": [30, 40, 20], 
		"prices": [
			[2, 3, 2, 4],
			[3, 2, 5, 1],
			[4, 3, 2, 6]
		]
	}`)

	message := Message{}

	if err := json.Unmarshal(msg, &message); err != nil {
		panic(err)
	}

	priceMatrix := initPriceMatrix(message)
	findBasicSolution(priceMatrix, &message) // TODO add copy of message
	printPM(*priceMatrix)

	fmt.Println("done")
}
