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
[First idea]

Take first 12 elements
Foreach next num in remainder

	Shutdown on of the subset of 12 elements
	Create a new subset of 12 appending next num
	Check if new subset value is bigger
*/
func Pick12Batteries(bank string) int {

	// // Convert bank to int[]
	var batteries []int = BankStrToIntArray(bank)
	var result []int = make([]int, 0, 12)

	remaining := 12

	for remaining > 0 {

		leftEndIdx := len(batteries) - remaining + 1
		if leftEndIdx == 0 {
			// Nothing to choose, append batteries to result
			result = append(result, batteries...)
			remaining = 0
			break
		}

		leftSide := batteries[0:leftEndIdx]
		rightSide := batteries[leftEndIdx:]

		max := slices.Max(leftSide)
		result = append(result, max)

		remaining--
		batteries = rightSide
	}

	return IntArrayToJoltageValue(result)

	// var subsetOf12 []int = batteries[0:12]
	// var remainder []int = batteries[12:]
	// valueOfSubset := IntArrayToJoltageValue(subsetOf12)

	// var bestSubsetOf12 []int = subsetOf12
	// var bestJoltage = valueOfSubset

	// for _, nextNum := range remainder {

	// 	for subsetIdx := range 12 {

	// 		// Exclude one of the original 12 numbers
	// 		var newSubsetOf12 []int = SliceExclude(subsetOf12, subsetIdx)
	// 		// Append next value
	// 		newSubsetOf12 = append(newSubsetOf12, nextNum)
	// 		// Recompute new joltage
	// 		newJoltage := IntArrayToJoltageValue(newSubsetOf12)

	// 		if newJoltage > bestJoltage {
	// 			bestSubsetOf12 = newSubsetOf12
	// 			bestJoltage = newJoltage
	// 		}
	// 	}

	// 	subsetOf12 = bestSubsetOf12
	// 	valueOfSubset = bestJoltage
	// }

	// return valueOfSubset
}

func Pick2Batteries(bank string) int {

	// Convert bank to int[]
	batteries := make([]int, 0, len(bank))

	for _, val := range bank {

		intVal, err := strconv.Atoi(string(val))
		if err != nil {
			panic("Error parsing char " + string(val))
		}

		batteries = append(batteries, intVal)
	}

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

	data := readData("test_input.txt")
	fmt.Printf("Part 2: %d\n", Part2(data))
	// fmt.Printf("%v\n", total)

	// fmt.Println(Pick12Batteries("234234234234278"))

}
