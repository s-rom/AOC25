package main

import (
	"bufio"
	"fmt"
	"strings"

	"log"
	"os"
)

func ReadData(filePath string) []string {
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

func ParseGrid(filePath string) *Grid[string] {
	data := ReadData(filePath)

	rows := len(data)
	columns := len(data[0])
	var pos Position
	grid := NewGrid[string](rows, columns)

	for row, line := range data {
		tokens := strings.Split(line, "")

		for col, tok := range tokens {
			pos.Row = row
			pos.Column = col
			grid.SetAt(pos, tok)
		}
	}

	return grid
}

func PropagateBreadth(grid *Grid[string], queue []Position) (uint64, []Position) {

	// Empty queue
	if len(queue) == 0 {
		return 0, queue
	}

	splits := uint64(0)

	currentPos := (queue)[0]
	queue = (queue)[1:] //pop

	// End of map
	if !grid.IsValidPosition(currentPos) {
		return 0, queue
	}

	posValuePtr := grid.GetPtrAt(currentPos)
	switch *posValuePtr {
	case ".":
		*posValuePtr = "|"
		nextPos := MoveDown(currentPos)
		queue = append(queue, nextPos)

	case "^":
		splits++
		leftPos := MoveLeft(currentPos)
		rightPos := MoveRight(currentPos)
		queue = append(queue, leftPos)
		queue = append(queue, rightPos)
	}

	return splits, queue
}

func FindStartPosition(grid *Grid[string]) Position {
	var startPosition Position
	for col := range grid.Columns {
		pos := Position{Row: 0, Column: col}
		if *grid.GetPtrAt(pos) == "S" {
			startPosition = pos
			break
		}
	}
	return startPosition
}

func Part1(grid *Grid[string]) uint64 {
	var total uint64 = 0

	queue := make([]Position, 0, 2*grid.Rows)
	queue = append(queue, MoveDown(FindStartPosition(grid)))

	var splits uint64
	for len(queue) != 0 {
		splits, queue = PropagateBreadth(grid, queue)
		total += splits
	}
	return total
}

func CountPossiblePaths(grid *Grid[string], cache map[Position]uint64, from Position) uint64 {
	paths := uint64(0)

	// At last row
	if grid.Rows-1 == from.Row {
		return 1
	}

	// downPosition := MoveDown(from)
	if !grid.IsValidPosition(from) {
		return 0
	}

	if val, ok := cache[from]; ok {
		// fmt.Print("Cache hit\n")
		return val
	}

	ptr := grid.GetPtrAt(from)

	switch *ptr {
	case ".":
		paths += CountPossiblePaths(grid, cache, MoveDown(from))
	case "^":
		paths += CountPossiblePaths(grid, cache, MoveLeft(from))
		paths += CountPossiblePaths(grid, cache, MoveRight(from))
	}

	cache[from] = paths

	return paths
}

type Path []Position

func Part2(grid *Grid[string]) uint64 {
	total := uint64(0)
	cache := make(map[Position]uint64)
	total = CountPossiblePaths(grid, cache, MoveDown(FindStartPosition(grid)))
	return total
}

func main() {
	filePath := "input_test.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	// rawData := readData(filePath)
	data := ParseGrid(filePath)
	data2 := ParseGrid(filePath)
	// data := ParseData(readData(filePath))
	var part1 uint64 = Part1(data)
	var part2 uint64 = Part2(data2)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

	StartViz(data)

}
