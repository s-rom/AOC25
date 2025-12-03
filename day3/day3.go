package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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

func BankStrToIntArray(bank string) []int {
	batteries := make([]int, 0, len(bank))
	for _, val := range bank {

		intVal, err := strconv.Atoi(string(val))
		if err != nil {
			panic("Error parsing char " + string(val))
		}

		batteries = append(batteries, intVal)
	}
	return batteries
}

func IntArrayToJoltageValue(batteries []int) int {

	str := ""
	for _, val := range batteries {
		str += strconv.Itoa(val)
	}

	value, _ := strconv.Atoi(str)
	return value
}

func SliceExclude[T any](s []T, idx int) []T {
	return append(s[0:idx], s[idx+1:]...)
}

/*
[Main idea]

Split the array into two sides: left and right
Find the max in left, right is the slice starting from maxIdx + 1
Left is sized so there are at least REM-1 batteries at the right
Repeat until result has 12 batteries

i.e.

If result is empty, remaining = 12

	=> Left = batteries[0:4] because there are AT LEAST 11 (REM - 1) in Right
		Take max,maxIdx of Left
		Right = [maxIdx + 1]+
		Right becomes Batteries
*/
func Pick12Batteries(bank string) int {

	// // Convert bank to int[]
	var batteries []int = BankStrToIntArray(bank)
	var result []int = make([]int, 0, 12)

	remaining := 12

	for remaining > 0 {

		leftEndIdx := len(batteries) - remaining + 1
		if leftEndIdx == 0 {
			// Left is empty, append batteries to result
			result = append(result, batteries...)
			remaining = 0
			break
		}

		leftSide := batteries[0:leftEndIdx]

		max := slices.Max(leftSide)
		maxIdx := slices.Index(leftSide, max)
		result = append(result, max)
		remaining--

		batteries = batteries[maxIdx+1:]
	}

	return IntArrayToJoltageValue(result)
}

func Pick2Batteries(bank string) int {

	// Convert bank to int[]
	var batteries []int = BankStrToIntArray(bank)

	var max2 int
	max := slices.Max(batteries[0 : len(batteries)-1])
	idxOfMax := slices.Index(batteries, max)

	if idxOfMax < len(batteries)-1 {
		// Max element is not the last.
		// Find max to the right of max
		max2 = slices.Max(batteries[idxOfMax+1:])

	} else {
		// Max element is last --> max2 is the last element
		max2 = batteries[len(batteries)-1]
	}

	joltage, _ := strconv.Atoi(strconv.Itoa(max) + strconv.Itoa(max2))

	return joltage
}

func Part1(data []string) int64 {
	var total int64 = 0
	for _, bank := range data {
		joltage := Pick2Batteries(bank)
		fmt.Printf("%s --> %d\n", bank, joltage)
		total += int64(joltage)
	}
	return total
}

func Part2(data []string) int64 {
	var total int64 = 0
	for _, bank := range data {
		joltage := Pick12Batteries(bank)
		fmt.Printf("%s --> %d\n", bank, joltage)
		total += int64(joltage)
	}
	return total
}

func main() {

	data := readData("input.txt")
	fmt.Printf("Part 1: %d\n", Part1(data))
	fmt.Printf("Part 2: %d\n", Part2(data))
}
