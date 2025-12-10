package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var pattern = regexp.MustCompile(`\[(.*)\] (.*) \{(.*)\}`)

const ON = true
const OFF = false

type LightSet uint32
type Button uint32 // Toggle mask

type MachineSchema struct {
	Goal                LightSet
	ButtonMasks         []Button
	JoltageRequirements []int
}

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

func NewLightSet() LightSet {
	return LightSet(0)
}

func PressButton(set LightSet, button Button) LightSet {
	return set ^ LightSet(button)
}

func IsGoalAchieved(m *MachineSchema, pressSequence []int) bool {
	testSet := LightSet(0)
	for _, buttonIndex := range pressSequence {
		button := m.ButtonMasks[buttonIndex]
		testSet = PressButton(testSet, button)
	}

	return m.Goal == testSet
}

func DeepCopySlice(other []int) []int {
	newSlice := make([]int, 0, len(other))
	newSlice = append(newSlice, other...)
	return newSlice
}

func ParseLightSet(str string) LightSet {
	set := LightSet(0)

	for idx, s := range str {
		if s == '#' {
			set = set | 1
		}

		if idx < len(str)-1 {
			set = set << 1
		}
	}
	return set
}

func ParseButtons(s string, totalLights int) []Button {
	buttons := make([]Button, 0, 10)

	tokens := strings.SplitSeq(s, " ")
	for token := range tokens {
		button := Button(0)
		// remove first and last => ( and )
		token = token[1 : len(token)-1]
		for numStr := range strings.SplitSeq(token, ",") {
			num, _ := strconv.Atoi(numStr)
			unsigned := uint32(1) << uint32(totalLights-1-num)
			button = button | Button(unsigned)
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
			ButtonMasks:         ParseButtons(match[2], len(match[1])),
			JoltageRequirements: ParseJoltage(match[3]),
		}
		machines = append(machines, machine)
	}
	return machines
}

// func key(slice []int) string {
// 	slices.Sort(slice)
// 	var builder strings.Builder
// 	for _, n := range slice {
// 		str := strconv.Itoa(n)
// 		builder.WriteString(str)
// 	}
// 	return builder.String()
// }

func MinNumberOfPresses(m *MachineSchema) int {

	minNumber := math.MaxInt
	queue := make([][]int, 0, 10)
	// Enqueue one sequence for each button available.
	// Test just one button press first
	for buttonIndex := range m.ButtonMasks {
		buttonSequence := make([]int, 0, 10)
		buttonSequence = append(buttonSequence, buttonIndex)
		queue = append(queue, buttonSequence)
	}

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:] // pop

		// If the current button sequence is bigger than the minium, skip it
		if len(front) > minNumber {
			continue
		}

		if IsGoalAchieved(m, front) {
			presses := len(front)
			if presses < minNumber {
				minNumber = presses
			}
		}

		for buttonIndex := range m.ButtonMasks {

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
	// m1 := machines[0]
	// set := LightSet(0)
	// set = PressButton(set, m1.ButtonMasks[0])
	// set = PressButton(set, m1.ButtonMasks[1])
	// set = PressButton(set, m1.ButtonMasks[2])

	// fmt.Println(strconv.FormatInt(int64(set), 2))

	// for _, m := range machines {
	// 	slices.Sort(m.Buttons)
	// }

	var part1 uint64 = Part1(machines)
	var part2 uint64 = Part2(rawData)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)
}
