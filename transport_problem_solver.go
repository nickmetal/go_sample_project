package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Cell struct {
	x        int     // row index
	y        int     // column index
	price    float64 // transporting price from producers to consumers
	consumed float64 // represent value of needs which was provided from consumer
	marked   bool    // mark as processed by NorthWest Method
}

func (c Cell) String() string {
	return fmt.Sprintf("s{x:%d y:%d con:%f mark:%v}", c.x, c.y, c.consumed, c.marked)
}

func (c Cell) Print() { fmt.Printf("%s\n", c) }
func printNeedsAndSourcesState(m Message, extra ...interface{}) {
	var s string
	if len(extra) == 1 {
		s = fmt.Sprintf("Extra: %v", extra[0])
	}

	fmt.Printf(
		"%vneeds: %v\nsources: %v\n",
		s,
		m.ConsumersNeeds,
		m.ProducersSources,
	)
}

type Message struct {
	Prices           [][]float64 `json:"prices"`
	ProducersSources []float64   `json:"producers_sources"`
	ConsumersNeeds   []float64   `json:"consumers_needs"`
}

type PriceMatrix [][]Cell

func IntContainsInSlice(value int, s []int) bool {
	for _, v := range s {
		if value == v {
			return true
		}
	}
	return false
}

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
			cell := Cell{x: rowIdx, y: columnIdx, price: cellValue}
			pm[rowIdx] = append(pm[rowIdx], cell)
		}
	}
	return &pm
}

func findBasicSolution(pm *PriceMatrix, m *Message) {

	for rId, row := range m.Prices {
		// for columnIdx, cellValue := range row {
		for cId, _ := range row {
			matrix := *(pm)
			inputs := *(m)

			if matrix[rId][cId].marked {
				fmt.Printf("  [value] for matrix[x,y]: %s IS SKIPPED =0 \n", matrix[rId][cId])
				continue
			}
			fmt.Printf("  [value] for matrix[x,y]: %s \n", matrix[rId][cId])

			if inputs.ProducersSources[rId] <= 0 {
				item := matrix[rId][cId]
				panic(fmt.Sprintf("ProducersSources is empty for (%v %v): %v\n", item.x, item.y, inputs.ProducersSources[rId]))
			}

			consumed := math.Min(inputs.ProducersSources[rId], inputs.ConsumersNeeds[cId])
			matrix[rId][cId].consumed = consumed

			printNeedsAndSourcesState(inputs)
			inputs.ConsumersNeeds[cId] -= consumed
			inputs.ProducersSources[rId] -= consumed
			printNeedsAndSourcesState(inputs)

			// if we provide full needs
			if inputs.ConsumersNeeds[cId] == 0 {
				fmt.Printf("found full: %s \n", matrix[rId][cId])
				markColumnCellsAsSkipped(pm, cId, rId)
				fmt.Println("")
				break
			}
			fmt.Println("")
		}
		fmt.Println(strings.Repeat("===", 15))
	}
}

// mark all cells in column except consumed
func markColumnCellsAsSkipped(pm *PriceMatrix, columnIndex int, rowIdx int) {
	matrix := *(pm)
	for rId := range matrix {
		matrix[rId][columnIndex].marked = true
		fmt.Printf("marked item: %s\n", matrix[rId][columnIndex])
	}
}

func printPM(pm PriceMatrix) {
	fmt.Println("priceMatrix")

	for _, row := range pm {
		rows := []string{}
		for _, col := range row {
			rows = append(rows, fmt.Sprintf("%s\t", col))
		}
		fmt.Println(strings.Join(rows, ""))
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
	fmt.Printf("sources in the end: %v", message.ProducersSources)

	fmt.Println("done")
}
