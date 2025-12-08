package main

import (
	"bufio"
	"log"
	"os"
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

func Part1(data []string) uint64 {
	total := uint64(0)

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
	var part1 uint64 = Part1(rawData)
	var part2 uint64 = Part2(rawData)

	// fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
