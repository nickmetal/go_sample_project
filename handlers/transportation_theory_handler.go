package handlers

// package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

type Cell struct {
	x              int     // row index
	y              int     // column index
	price          float64 // transporting price from producers to consumers
	consumed       float64 // represent value of needs which was provided from consumer
	cellDifference float64 // ΔCij = Cij – (Ui + Vj ) difference
}

func (c Cell) String() string {
	return fmt.Sprintf("s{x:%d y:%d con:%g price:%g diff:%g }", c.x, c.y, c.consumed, c.price, c.cellDifference)
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

func (pm *PriceMatrix) findBasicSolution(inputs Message) error {
	log.Println("findBasicSolution")

	matrix := *(pm)

	for cId, column := range matrix {
		for rId, cell := range column {

			if cell.consumed != 0 {
				continue
			}
			// fmt.Printf("  [value] for matrix[x,y]: %s \n", cell)

			if inputs.ProducersSources[rId] < 0 {
				item := matrix[rId][cId]
				return fmt.Errorf("ProducersSources is empty for (%v %v): %v\n", item.x, item.y, inputs.ProducersSources[rId])
			}

			consumed := math.Min(inputs.ProducersSources[rId], inputs.ConsumersNeeds[cId])
			cell.consumed = consumed
			matrix[cId][rId] = cell

			if inputs.ConsumersNeeds[cId] == 0 && rId == 0 {
				return errors.New(fmt.Sprintf("found cell with out needs: %s \n", cell))
			}

			inputs.ConsumersNeeds[cId] -= consumed
			inputs.ProducersSources[rId] -= consumed
		}
	}
	return nil
}

func (m PriceMatrix) validateBasicSolution(message Message) error {
	log.Println("validateBasicSolution")

	// Проверка, что все поставщики израсходовали свои запасы
	sourcesSum := 0.0
	for _, value := range message.ProducersSources {
		sourcesSum += value
	}

	if sourcesSum != 0 {
		return fmt.Errorf("All sources should  be empty. Now total sum of sources: %g", sourcesSum)
	}

	// Проверка, что все потребители получили желаемое количество единиц товара
	needsSum := 0.0
	for _, value := range message.ConsumersNeeds {
		needsSum += value
	}

	if needsSum != 0 {
		return fmt.Errorf("All need should be empty. Now total sum of needs: %g", needsSum)
	}

	/*
		Проверка плана на вырожденность.
		Базисных ячеек таблицы должно быть не менее m+n-1
		где m и n — соответственно, число поставщиков и потребителей,
		иначе решение считается вырожденным
	*/
	basicCellNum := 0
	for _, column := range m {
		for _, cell := range column {
			if cell.consumed != 0 {
				basicCellNum += 1
			}
		}
	}

	if basicCellNum < len(m)+len(m[0])-1 {
		return fmt.Errorf("basic plan вырожденный: basicCellNum: %d, m+n=%d", basicCellNum, len(m)+len(m[0])-1)
	}
	return nil
}

func (m *PriceMatrix) calculateDifferencesForOptimum(sp *[]float64, cp *[]float64) bool {
	log.Println("calculateDifferencesForOptimum")

	sourcesPotentials := *sp
	consumerPotentials := *cp
	foundOptimum := true

	for cId, column := range *m {
		for rId, cell := range column {
			// skip consumed cells
			if cell.consumed != 0 {
				continue
			}
			cellDifference := cell.price - sourcesPotentials[rId] - consumerPotentials[cId]
			cell.cellDifference = cellDifference
			(*m)[cId][rId] = cell

			if cellDifference < 0 {
				foundOptimum = false
			}
		}
	}
	return foundOptimum
}

func (m *PriceMatrix) calculatePotentials(message Message) (*[]float64, *[]float64) {
	matrix := *(m)
	consumerPotentials := make([]float64, len(message.ConsumersNeeds))
	// todo rename to providers etc
	sourcesPotentials := make([]float64, len(message.ProducersSources))

	var currentConsumerP float64
	var currentSourceP float64

	for rId, row := range message.Prices {
		for cId, price := range row {
			if cId == 0 && rId == 0 {
				consumerPotentials[0] = price
				continue
			}

			// non basic cell
			if matrix[cId][rId].consumed == 0 {
				continue
			}
			if rId != 0 && sourcesPotentials[rId] == 0 {
				currentSourceP = price - currentConsumerP
				sourcesPotentials[rId] = currentSourceP
			}

			if consumerPotentials[cId] == 0 {
				currentConsumerP = price - currentSourceP
				consumerPotentials[cId] = currentConsumerP
			}
		}
	}
	return &sourcesPotentials, &consumerPotentials
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

func initPriceMatrix(m Message) *PriceMatrix {
	log.Println("initPriceMatrix")

	pm := make(PriceMatrix, 0)

	for y, _ := range m.ConsumersNeeds {
		columnSlice := make([]Cell, 0)

		for x, _ := range m.ProducersSources {
			columnSlice = append(columnSlice, Cell{x: y, y: x, price: m.Prices[x][y]})
		}
		pm = append(pm, columnSlice)
	}
	return &pm
}

func printPM(pm PriceMatrix) {
	log.Println("priceMatrix")

	for _, row := range pm {
		rows := []string{}
		for _, col := range row {
			rows = append(rows, fmt.Sprintf("%s\t", col))
		}
		log.Println(strings.Join(rows, ""))
	}
}

func Solve(message Message) error {
	basicSolutionMatrix := initPriceMatrix(message)
	var err error
	err = basicSolutionMatrix.findBasicSolution(message)
	if err != nil {
		return err
	}
	err = basicSolutionMatrix.validateBasicSolution(message)
	if err != nil {
		return err
	}

	sourcesPotentials, consumerPotentials := basicSolutionMatrix.calculatePotentials(message)

	foundOptimum := basicSolutionMatrix.calculateDifferencesForOptimum(
		sourcesPotentials,
		consumerPotentials,
	)
	if !foundOptimum {
		return errors.New("unable to find optimum for input data")
	}
	printPM(*basicSolutionMatrix)
	return nil
}

func TransportIssueHandler(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		body, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		if err != nil {
			http.Error(resp, err.Error(), 500)
			log.Println("POST body read error", err)
			return
		}
		log.Println("POST request", string(body))

		message := Message{}

		if err := json.Unmarshal(body, &message); err != nil {
			log.Println("POST body unmarshal error", err)
			http.Error(resp, err.Error(), 500)
			return
		}

		err = Solve(message)
		if err == nil {
			resp.Write([]byte("{\"data\": {}}"))
		} else {
			http.Error(resp, err.Error(), 400)
		}

	default:
		http.Error(resp, "unsupported method", 400)
		log.Printf("unsupported method: %s\n", req.Method)
	}
}

// func main() {
// 	msg := []byte(`{
// 		"consumers_needs": [20, 30, 30, 10],
// 		"producers_sources": [30, 40, 20],
// 		"prices": [
// 			[2, 3, 2, 4],
// 			[3, 2, 5, 1],
// 			[4, 3, 2, 6]
// 		]
// 	}`)

// 	message := Message{}

// 	if err := json.Unmarshal(msg, &message); err != nil {
// 		panic(err)
// 	}
// 	Solve(message)
// }
