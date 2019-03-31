package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Cell struct {
	x int // row index
	y int // column index
	// price    float64 // transporting price from producers to consumers
	consumed float64 // represent value of needs which was provided from consumer
	marked   bool    // mark as processed by NorthWest Method
}

func (c Cell) String() string {
	return fmt.Sprintf("s{x:%d y:%d con:%g mark:%v}", c.x, c.y, c.consumed, c.marked)
}

func (c Cell) Print() { fmt.Printf("%s\n", c) }
func printNeedsAndSourcesState(m Message, extra ...interface{}) {
	var s string
	if len(extra) == 1 {
		s = fmt.Sprintf("Extra: %v", extra[0])
	}

	fmt.Printf(
		"%vneeds:   %v\nsources: %v\n",
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

	column based matrix: slice of matrix columns
*/
func initPriceMatrix(m Message) *PriceMatrix {
	// consumersNums := len(m.ConsumersNeeds)
	pm := make(PriceMatrix, 0)

	for y, _ := range m.ConsumersNeeds {
		columnSlice := make([]Cell, 0)

		for x, _ := range m.ProducersSources {
			columnSlice = append(columnSlice, Cell{x: y, y: x}) // todo add transportation price
		}
		pm = append(pm, columnSlice)
	}
	return &pm
}

func findBasicSolution(pm *PriceMatrix, m *Message) {
	matrix := *(pm)
	inputs := *(m)
	// producersCount := len(inputs.ProducersSources) consumerXYneeds
	// consumersCount := len(inputs.ConsumersNeeds) producerXYsources

	for cId, column := range matrix {
		for rId, cell := range column {

			if cell.consumed != float64(0) {
				fmt.Printf("  [value] for matrix[x,y]: %s IS SKIPPED =0 \n", cell)
				continue
			}
			fmt.Printf("  [value] for matrix[x,y]: %s \n", cell)

			if inputs.ProducersSources[rId] < 0 {
				item := matrix[rId][cId]
				panic(fmt.Sprintf("ProducersSources is empty for (%v %v): %v\n", item.x, item.y, inputs.ProducersSources[rId]))
			}

			consumed := math.Min(inputs.ProducersSources[rId], inputs.ConsumersNeeds[cId])
			cell.consumed = consumed
			cell.marked = true
			matrix[cId][rId] = cell

			if inputs.ConsumersNeeds[cId] == 0 && rId == 0 {
				panic(fmt.Sprintf("found cell with out needs: %s \n", cell))
			}

			printNeedsAndSourcesState(inputs)
			inputs.ConsumersNeeds[cId] -= consumed
			inputs.ProducersSources[rId] -= consumed
			printNeedsAndSourcesState(inputs)

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

	basicSolutionMatrix := initPriceMatrix(message)
	// priceMatrix = findBasicSolution(priceMatrix, &message) // TODO add copy of message
	findBasicSolution(basicSolutionMatrix, &message) // TODO add copy of message
	errors := validateBasicSolution(basicSolutionMatrix)
	if errors != nil {
		panic(errors)
	}

	printPM(*priceMatrix)
	fmt.Printf("sources in the end: %v", message.ProducersSources)

	fmt.Println("done")
}
