package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
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
	ButtonValues        [][]int
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

func ParseButtons(s string, totalLights int) ([]Button, [][]int) {
	buttonMasks := make([]Button, 0, 10)
	buttonValues := make([][]int, 0, 10)

	tokens := strings.SplitSeq(s, " ")
	for token := range tokens {
		mask := Button(0)
		// remove first and last => ( and )
		token = token[1 : len(token)-1]

		values := make([]int, 0, totalLights)

		for numStr := range strings.SplitSeq(token, ",") {
			num, _ := strconv.Atoi(numStr)
			unsigned := uint32(1) << uint32(totalLights-1-num)
			mask = mask | Button(unsigned)
			values = append(values, num)
		}
		buttonValues = append(buttonValues, values)
		buttonMasks = append(buttonMasks, mask)
	}

	return buttonMasks, buttonValues
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

		masks, buttonValues := ParseButtons(match[2], len(match[1]))

		machine := MachineSchema{
			Goal:                ParseLightSet(match[1]),
			ButtonMasks:         masks,
			ButtonValues:        buttonValues,
			JoltageRequirements: ParseJoltage(match[3]),
		}
		machines = append(machines, machine)
	}
	return machines
}

func MinNumberOfPressesLightGoal(m *MachineSchema) int {

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
			// ==> unnecessary
			if buttonIndex == front[len(front)-1] {
				continue
			}

			// Avoid chosing a button already pushed
			idx := slices.Index(front, buttonIndex)
			if idx != -1 {
				continue
			}

			newSequence := DeepCopySlice(front)
			newSequence = append(newSequence, buttonIndex)
			queue = append(queue, newSequence)
		}
	}

	return minNumber
}

func SliceEquals(s1 []int, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func GetJoltageValue(m *MachineSchema, buttonSeq []int) []int {

	result := make([]int, len(m.JoltageRequirements))

	for _, buttonIdx := range buttonSeq {
		for _, joltageIndex := range m.ButtonValues[buttonIdx] {
			result[joltageIndex]++
		}
	}

	return result
}

func Part1(data []MachineSchema) uint64 {
	total := uint64(0)
	progress := 0
	for _, m := range data {

		n := MinNumberOfPressesLightGoal(&m)
		total += uint64(n)

		progress++
		// fmt.Println(float64(progress) / float64(len(data)) * 100.0)
	}
	return total
}

type EquationSet struct {
	Equations    []Equation
	OperandNames []string
}

type Equation struct {
	Result   int
	Operands []string
}

func GetEquations(m MachineSchema) EquationSet {
	operandNames := make([]string, 0, len(m.JoltageRequirements))

	for opIdx := range len(m.ButtonValues) {
		// 'a' = 97
		operandNames = append(operandNames, string(rune(97+opIdx)))
	}

	equations := make([]Equation, 0, 5)
	for joltageIdx, joltage := range m.JoltageRequirements {
		eq := Equation{
			Result:   joltage,
			Operands: make([]string, 0, len(m.JoltageRequirements)),
		}

		for buttonIndex, button := range m.ButtonValues {
			idx := slices.Index(button, joltageIdx)
			if idx != -1 { // Found
				buttonName := operandNames[buttonIndex]
				eq.Operands = append(eq.Operands, buttonName)
			}
		}

		equations = append(equations, eq)
	}

	return EquationSet{
		Equations:    equations,
		OperandNames: operandNames,
	}
}

func GenerateSMT2Code(eqs EquationSet) string {
	var builder strings.Builder

	builder.WriteString("(push 1)\n")

	for _, opName := range eqs.OperandNames {
		builder.WriteString("(declare-const ")
		builder.WriteString(opName)
		builder.WriteString(" Int)\n")
	}

	for _, opName := range eqs.OperandNames {
		builder.WriteString("(assert (>= ")
		builder.WriteString(opName)
		builder.WriteString(" 0))\n")
	}

	for _, eq := range eqs.Equations {
		builder.WriteString("(assert (= (+ ")
		for _, opName := range eq.Operands {
			builder.WriteString(opName)
			builder.WriteString(" ")
		}
		builder.WriteString(") ")
		builder.WriteString(strconv.Itoa(eq.Result))
		builder.WriteString("))\n")
	}

	for _, opName := range eqs.OperandNames {
		builder.WriteString("(minimize ")
		builder.WriteString(opName)
		builder.WriteString(")\n")
	}

	builder.WriteString("(check-sat)\n(get-value (")

	for _, opName := range eqs.OperandNames {
		builder.WriteString(opName)
		builder.WriteString(" ")
	}

	builder.WriteString("))\n")

	builder.WriteString("(pop 1)\n")

	return builder.String()
}

func Part2(data []MachineSchema) uint64 {
	total := uint64(0)

	// TODO: Solve using z3 library

	return total
}

func main() {

	filePath := "input_test.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	rawData := readData(filePath)
	machines := ParseData(rawData)

	// var part1 uint64 = Part1(machines)
	// var part2 uint64 = Part2(machines)

	// fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

	for _, m := range machines {
		eq := GetEquations(m)
		fmt.Println(GenerateSMT2Code(eq))
	}

	// fmt.Println(eq)

}
