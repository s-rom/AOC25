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

type IngredientRange struct {
	Start uint64
	End   uint64
}

func (irange IngredientRange) IsInRange(value uint64) bool {
	return value >= irange.Start && value <= irange.End
}

func (irange IngredientRange) RangeLength() uint64 {
	return irange.End - irange.Start + 1
}

func IntersectRange(a IngredientRange, b IngredientRange) IngredientRange {
	return IngredientRange{
		Start: max(a.Start, b.Start),
		End:   min(a.End, b.End),
	}
}

func TestIntersection(a IngredientRange, b IngredientRange) bool {
	return a.IsInRange(b.Start) || a.IsInRange(b.End) || b.IsInRange(a.Start) || b.IsInRange(a.End)
}

type IngredientData struct {
	Ranges []IngredientRange
	Ids    []uint64
}

func ParseData(data []string) IngredientData {

	result := IngredientData{
		Ranges: make([]IngredientRange, 0, 1024),
		Ids:    make([]uint64, 0, 1024),
	}

	// Parse ranges
	idx := 0
	for lineIdx, line := range data {
		if line == "" {
			idx = lineIdx + 1
			break
		}

		tokens := strings.Split(line, "-")
		start, _ := strconv.ParseUint(tokens[0], 10, 64)
		end, _ := strconv.ParseUint(tokens[1], 10, 64)

		result.Ranges = append(result.Ranges, IngredientRange{
			Start: start,
			End:   end,
		})
	}

	for i := idx; i < len(data); i++ {
		id, _ := strconv.ParseUint(data[i], 10, 64)
		result.Ids = append(result.Ids, id)
	}

	return result
}

func Part1(data IngredientData) uint64 {
	var total uint64 = 0

	for _, id := range data.Ids {

		fresh := false

		for _, ingredientRange := range data.Ranges {
			if ingredientRange.IsInRange(id) {
				fresh = true
				break
			}
		}

		if fresh {
			total++
		}
	}

	return total
}

func Part2(data IngredientData) uint64 {

	total := uint64(0)
	for idxA, rangeA := range data.Ranges {

		total += rangeA.RangeLength()

		for idxB := idxA + 1; idxB < len(data.Ranges); idxB++ {
			rangeB := &data.Ranges[idxB]

			if TestIntersection(rangeA, *rangeB) {
				intersection := IntersectRange(rangeA, *rangeB)
				total -= intersection.RangeLength()
			}
		}
	}

	return total
}

func main() {
	filePath := "input.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	data := ParseData(readData(filePath))
	var part1 uint64 = Part1(data)
	var part2 uint64 = Part2(data)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
