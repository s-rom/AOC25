package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var pattern = regexp.MustCompile(`\[(.*)\] (.*) \{(.*)\}`)

const ON = true
const OFF = false

type LightSet []bool
type Button []int

type MachineSchema struct {
	Goal                LightSet
	Buttons             []Button
	JoltageRequirements []int
}

// func (m MachineSchema) String() string {
// 	var builder strings.Builder

// 	builder.WriteRune('[')
// 	for _, light := range m.Goal {
// 		if light {
// 			builder.WriteRune('#')
// 		} else {
// 			builder.WriteRune('.')
// 		}
// 	}
// 	builder.WriteString("]\n")

// 	return builder.String()
// }

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

func ParseLightSet(s string) LightSet {
	set := make(LightSet, 0, len(s))

	for _, s := range s {
		var value bool
		if s == '.' {
			value = OFF
		} else {
			value = ON
		}
		set = append(set, value)
	}
	return set
}

func ParseButons(s string) []Button {
	buttons := make([]Button, 0, 10)

	tokens := strings.SplitSeq(s, " ")
	for token := range tokens {
		button := make(Button, 0, 10)
		// remove first and last => ( and )
		token = token[1 : len(token)-2]
		for numStr := range strings.SplitSeq(token, ",") {
			num, _ := strconv.Atoi(numStr)
			button = append(button, num)
		}
		buttons = append(buttons, button)
	}

	return buttons
}

func ParseJoltage(s string) []int {
	joltage := make([]int, 0, 10)
	for strNum := range strings.SplitSeq(s, ",") {
		num, _ := strconv.Atoi(strNum)
		joltage = append(joltage, num)
	}
	return joltage
}

func ParseData(data []string) {
	for _, line := range data {
		match := pattern.FindStringSubmatch(line)

		machine := MachineSchema{
			Goal:                ParseLightSet(match[1]),
			Buttons:             ParseButons(match[2]),
			JoltageRequirements: ParseJoltage(match[3]),
		}

		fmt.Println(machine)
	}
}

func MinNumberOfPresses(m *MachineSchema) int {
	numberOfPresses := 0

	return numberOfPresses
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
	ParseData(rawData)
	var part1 uint64 = Part1(rawData)
	var part2 uint64 = Part2(rawData)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
