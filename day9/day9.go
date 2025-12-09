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

type IVec2 struct {
	X int
	Y int
}

func ParseData(data []string) []IVec2 {
	result := make([]IVec2, 0, len(data))
	for _, line := range data {
		tokens := strings.Split(line, ",")
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])

		result = append(result, IVec2{
			X: x,
			Y: y,
		})
	}
	return result
}

func GetRectArea(a IVec2, b IVec2) uint64 {
	dx := max(a.X, b.X) - min(a.X, b.X) + 1
	dy := max(a.Y, b.Y) - min(a.Y, b.Y) + 1
	return uint64(dx * dy)
}

func Part1(data []IVec2) uint64 {
	total := uint64(0)

	for idxA, a := range data {
		for idxB := idxA + 1; idxB < len(data); idxB++ {
			b := data[idxB]
			area := GetRectArea(a, b)
			if area > total {
				total = area
			}
		}
	}

	return total
}

func Part2(data []string) uint64 {
	total := uint64(0)

	return total
}

func main() {
	filePath := "input_test.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	rawData := readData(filePath)
	data := ParseData(rawData)
	var part1 uint64 = Part1(data)
	var part2 uint64 = Part2(rawData)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
