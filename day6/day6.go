package main

import (
	"bufio"
	"fmt"

	"log"
	"os"
	"strconv"
	"strings"
)

func readData(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Could not open input file")
	}

	defer file.Close()

	var buffer []string = make([]string, 0, 200)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
	}

	return buffer
}

func ParseData(data []string) *Grid[string] {
	parsedData := []string{}

	rows := 0
	columns := 0

	for _, line := range data {
		columns = 0
		tokens := strings.SplitSeq(line, " ")

		for token := range tokens {
			if token != "" {
				parsedData = append(parsedData, token)
				columns++
			}
		}
		rows++
	}

	grid := NewGrid[string](rows, columns)
	grid.Data = parsedData
	return grid
}

func AddColumn(grid *Grid[string], column int) uint64 {
	total := uint64(0)

	for row := 0; row < grid.Rows-1; row++ {
		val := grid.GetPtrAt(Position{Row: row, Column: column})
		num, _ := strconv.ParseUint(*val, 10, 64)
		total += num
	}

	return total
}

func MultColumn(grid *Grid[string], column int) uint64 {
	total := uint64(1)

	for row := 0; row < grid.Rows-1; row++ {
		val := grid.GetPtrAt(Position{Row: row, Column: column})
		num, _ := strconv.ParseUint(*val, 10, 64)
		total *= num
	}

	return total
}

func Part1(grid *Grid[string]) uint64 {
	var total uint64 = 0

	for col := range grid.Columns {

		op := grid.GetPtrAt(Position{Row: grid.Rows - 1, Column: col})
		switch *op {
		case "*":
			total += MultColumn(grid, col)
		case "+":
			total += AddColumn(grid, col)
		}
	}

	return total
}

func Part2(rawData []string) uint64 {
	total := uint64(0)

	operations := make([]string, 0, 50)
	opTokens := strings.Split(rawData[len(rawData)-1], " ")
	for _, op := range opTokens {
		if op != "" {
			operations = append(operations, op)
		}
	}

	lineLength := len(rawData[0])

	operationIdx := 0
	var acc uint64
	switch operations[len(operations)-1] {
	case "*":
		acc = 1
	case "+":
		acc = 0
	}

	for col := lineLength - 1; col >= 0; col-- {
		currentOperation := operations[len(operations)-operationIdx-1]

		allSpaces := true
		for row := 0; row < len(rawData)-1; row++ {
			if rawData[row][col] != ' ' {
				allSpaces = false
			}
		}

		if allSpaces {
			operationIdx++ // next operation set
			total += acc

			nextOperation := operations[len(operations)-operationIdx-1]
			switch nextOperation {
			case "*":
				acc = 1
			case "+":
				acc = 0
			}
			continue
		}

		accNum := ""
		for row := 0; row < len(rawData)-1; row++ {
			if rawData[row][col] == ' ' {
				continue
			}
			accNum = string(accNum) + string(rawData[row][col])
		}

		num, _ := strconv.ParseUint(accNum, 10, 64)
		switch currentOperation {
		case "*":
			acc *= num
		case "+":
			acc += num
		}

	}
	total += acc
	return total
}

func main() {
	filePath := "input_test.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	rawData := readData(filePath)
	data := ParseData(rawData)
	// data := ParseData(readData(filePath))
	var part1 uint64 = Part1(data)
	var part2 uint64 = Part2(rawData)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
