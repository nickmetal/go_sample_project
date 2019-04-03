package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

type Cell struct {
	x        int     // row index
	y        int     // column index
	price    float64 // transporting price from producers to consumers
	consumed float64 // represent value of needs which was provided from consumer
	// potential float64
}

func (c Cell) String() string {
	return fmt.Sprintf("s{x:%d y:%d con:%g, price:%g}", c.x, c.y, c.consumed, c.price)
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

func (pm *PriceMatrix) findBasicSolution(inputs Message) {
	matrix := *(pm)

	for cId, column := range matrix {
		for rId, cell := range column {

			// todo check without float64 and remote it
			if cell.consumed != 0 {
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
			matrix[cId][rId] = cell

			if inputs.ConsumersNeeds[cId] == 0 && rId == 0 {
				panic(fmt.Sprintf("found cell with out needs: %s \n", cell))
			}

			// printNeedsAndSourcesState(inputs)
			inputs.ConsumersNeeds[cId] -= consumed
			inputs.ProducersSources[rId] -= consumed
			// printNeedsAndSourcesState(inputs)

			fmt.Println("")
		}
		fmt.Println(strings.Repeat("===", 15))
	}
}

func (m PriceMatrix) validateBasicSolution(message Message) []error {
	errorSlice := make([]error, 0)

	// Проверка, что все поставщики израсходовали свои запасы
	sourcesSum := 0.0
	for _, value := range message.ProducersSources {
		sourcesSum += value
	}

	if sourcesSum != 0 {
		errorSlice = append(
			errorSlice,
			fmt.Errorf("All sources should  be empty. Now total sum of sources: %g", sourcesSum),
		)
	}

	// Проверка, что все потребители получили желаемое количество единиц товара
	needsSum := 0.0
	for _, value := range message.ConsumersNeeds {
		needsSum += value
	}

	if needsSum != 0 {
		errorSlice = append(
			errorSlice,
			fmt.Errorf("All need should be empty. Now total sum of needs: %g", needsSum),
		)
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
		errorSlice = append(
			errorSlice,
			fmt.Errorf("basic plan вырожденный: basicCellNum: %d, m+n=%d", basicCellNum, len(m)+len(m[0])-1),
		)
	}

	if len(errorSlice) == 0 {
		return nil
	}
	return errorSlice
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
		fmt.Printf("[%d]: sourPots %v \n", rId, sourcesPotentials)
		fmt.Printf("[%d]: consPots %v \n", rId, consumerPotentials)
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
			columnSlice = append(columnSlice, Cell{x: y, y: x, price: m.Prices[x][y]})
		}
		pm = append(pm, columnSlice)
	}
	return &pm
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

func Solve(message Message) {
	// msg := []byte(`{
	// 	"consumers_needs": [20, 30, 30, 10],
	// 	"producers_sources": [30, 40, 20],
	// 	"prices": [
	// 		[2, 3, 2, 4],
	// 		[3, 2, 5, 1],
	// 		[4, 3, 2, 6]
	// 	]
	// }`)

	// message := Message{}

	// if err := json.Unmarshal(msg, &message); err != nil {
	// 	panic(err)
	// }

	basicSolutionMatrix := initPriceMatrix(message)
	basicSolutionMatrix.findBasicSolution(message)
	errorSlice := basicSolutionMatrix.validateBasicSolution(message)

	if errorSlice != nil {
		fmt.Println("Errors:")
		for _, validationError := range errorSlice {
			fmt.Println(validationError)
		}
		return
	}
	sourcesPotentials, consumerPotentials := basicSolutionMatrix.calculatePotentials(message)
	fmt.Printf("sourcesPotentials: %v", sourcesPotentials)
	fmt.Printf("consumerPotentials : %v", consumerPotentials)

	printPM(*basicSolutionMatrix)
	fmt.Printf("sources in the end: %v", message.ProducersSources)

	fmt.Println("done")
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

		Solve(message)
		resp.Write([]byte("{}"))

	default:
		http.Error(resp, "unsupported method", 400)
		log.Printf("unsupported method: %s\n", req.Method)
	}
}
