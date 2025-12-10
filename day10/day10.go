package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"sort"
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
// var builder strings.Builder

// builder.WriteRune('[')
// for _, light := range m.Goal {
// if light {
// builder.WriteRune('#')
// } else {
// builder.WriteRune('.')
// }
// }
// builder.WriteString("]\n")

// return builder.String()
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

func NewLightSet(size int) LightSet {
	set := make(LightSet, size)
	for i := range size {
		set[i] = OFF
	}
	return set
}

func IsGoalAchieved(m *MachineSchema, pressSequence []int) bool {
	testSet := NewLightSet(len(m.Goal))
	for _, buttonIndex := range pressSequence {
		for _, lightIdx := range m.Buttons[buttonIndex] {
			testSet[lightIdx] = !testSet[lightIdx] // toggle light
		}
	}

	for i := range m.Goal {
		if testSet[i] != m.Goal[i] {
			return false
		}
	}
	return true
}

func DeepCopySlice(other []int) []int {
	newSlice := make([]int, 0, len(other))
	newSlice = append(newSlice, other...)
	return newSlice
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

func ParseButtons(s string) []Button {
	buttons := make([]Button, 0, 10)

	tokens := strings.SplitSeq(s, " ")
	for token := range tokens {
		button := make(Button, 0, 10)
		// remove first and last => ( and )
		token = token[1 : len(token)-1]
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

func ParseData(data []string) []MachineSchema {
	machines := make([]MachineSchema, 0, len(data))

	for _, line := range data {
		match := pattern.FindStringSubmatch(line)

		machine := MachineSchema{
			Goal:                ParseLightSet(match[1]),
			Buttons:             ParseButtons(match[2]),
			JoltageRequirements: ParseJoltage(match[3]),
		}
		machines = append(machines, machine)
	}
	return machines
}

func key(slice []int) string {
	slices.Sort(slice)
	var builder strings.Builder
	for _, n := range slice {
		str := strconv.Itoa(n)
		builder.WriteString(str)
	}
	return builder.String()
}

func MinNumberOfPresses(m *MachineSchema) int {

	minNumber := math.MaxInt
	queue := make([][]int, 0, 10)
	cache := make(map[string]int)
	// Enqueue one sequence for each button available.
	// Test just one button press first
	for buttonIndex := range m.Buttons {
		buttonSequence := make([]int, 0, 10)
		buttonSequence = append(buttonSequence, buttonIndex)
		queue = append(queue, buttonSequence)
	}

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:] // pop
		key := key(front)

		if val, found := cache[key]; found {
			if val >= minNumber {
				continue
			}
		}

		// If the current button sequence is bigger than the minium, skip it
		if len(front) > minNumber {
			continue
		}

		if IsGoalAchieved(m, front) {
			presses := len(front)
			cache[key] = presses
			if presses < minNumber {
				minNumber = presses
			}
		}

		for buttonIndex := range m.Buttons {

			// Avoid pressing the last button again
			if buttonIndex == front[len(front)-1] {
				continue
			}

			newSequence := DeepCopySlice(front)
			newSequence = append(newSequence, buttonIndex)
			queue = append(queue, newSequence)
		}
	}

	return minNumber
}

func Part1(data []MachineSchema) uint64 {
	total := uint64(0)
	progress := 0
	for _, m := range data {

		n := MinNumberOfPresses(&m)
		total += uint64(n)

		progress++
		fmt.Println(float64(progress) / float64(len(data)) * 100.0)
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
	machines := ParseData(rawData)

	for _, m := range machines {
		sort.Slice(m.Buttons, func(i, j int) bool {
			return len(m.Buttons[i]) < len(m.Buttons[j])
		})
	}

	var part1 uint64 = Part1(machines)
	var part2 uint64 = Part2(rawData)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)
}
