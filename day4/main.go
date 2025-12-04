package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func ParseGrid(filePath string) *Grid[string] {
	data := readData(filePath)

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

func RemoveRolls(grid *Grid[string]) int {
	total := 0

	var pos Position
	toRemove := make([]Position, 20)

	for r := range grid.Rows {
		for c := range grid.Columns {
			pos.Row = r
			pos.Column = c

			value := grid.GetPtrAt(pos)
			if *value != "@" {
				continue
			}

			nearSum := 0
			for _, neigh := range grid.Get8Neighbours(pos) {
				if *grid.GetPtrAt(neigh) == "@" {
					nearSum++
				}
			}

			if nearSum < 4 {
				total++

				// Remove it
				toRemove = append(toRemove, pos)
			}
		}
	}

	for _, removePos := range toRemove {
		grid.SetAt(removePos, ".")
	}

	return total
}

func main() {
	filePath := "input.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	g := ParseGrid(filePath)

	part2 := 0
	removed := RemoveRolls(g)
	part1 := removed
	for removed != 0 {
		part2 += removed
		removed = RemoveRolls(g)
	}

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
