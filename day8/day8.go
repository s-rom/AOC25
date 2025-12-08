package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type IVec3 struct {
	X int
	Y int
	Z int
}

type Pair struct {
	Box1     IVec3
	Box2     IVec3
	Distance float64
}

func (a IVec3) Equals(b IVec3) bool {
	return a.X == b.X &&
		a.Y == b.Y &&
		a.Z == b.Z
}

func (a Pair) Equals(b Pair) bool {
	return a.Box1.Equals(b.Box1) || a.Box2.Equals(b.Box1) || a.Box1.Equals(b.Box2) || a.Box2.Equals(b.Box1)
}

func Distance(a IVec3, b IVec3) float64 {
	res := math.Pow(float64(a.X)-float64(b.X), 2) +
		math.Pow(float64(a.Y)-float64(b.Y), 2) +
		math.Pow(float64(a.Z)-float64(b.Z), 2)

	return math.Sqrt(res)
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

func ParseData(data []string) []IVec3 {
	result := make([]IVec3, len(data))
	for idx, line := range data {

		tokens := strings.Split(line, ",")

		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		z, _ := strconv.Atoi(tokens[2])

		result[idx] = IVec3{
			X: x,
			Y: y,
			Z: z,
		}
	}

	return result
}

func IsPairLinked(data *BoxesData, pair Pair) bool {
	a, foundA := data.CircuitsTable[pair.Box1]
	b, foundB := data.CircuitsTable[pair.Box2]

	return foundA && foundB && a == b
}

func GetSortedPairs(data *BoxesData) []Pair {
	pairs := make([]Pair, 0, len(data.Boxes)*len(data.Boxes))

	// fmt.Printf("Looking for next pair\n")

	for idxA, boxA := range data.Boxes {
		for idxB := idxA + 1; idxB < len(data.Boxes); idxB++ {
			boxB := data.Boxes[idxB]
			pairs = append(pairs,
				Pair{Box1: boxA, Box2: boxB, Distance: Distance(boxA, boxB)})
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Distance < pairs[j].Distance
	})

	return pairs
}

type BoxesData struct {
	Boxes         []IVec3
	CircuitsTable map[IVec3]int
	DistanceCache map[Pair]float64
	NextCircuitId int
}

func MergeCircuits(minPair Pair, data *BoxesData) {
	circuit1, box1InCircuit := data.CircuitsTable[minPair.Box1]
	circuit2, box2InCircuit := data.CircuitsTable[minPair.Box2]

	// fmt.Printf("Merging %v with %v\n", minPair.Box1, minPair.Box2)

	// None is in a circuit group
	if !box1InCircuit && !box2InCircuit {
		newCircuitId := data.NextCircuitId
		data.NextCircuitId++
		data.CircuitsTable[minPair.Box1] = newCircuitId
		data.CircuitsTable[minPair.Box2] = newCircuitId

		// fmt.Printf("\tNone in a citcuit, new id: %d\n", newCircuitId)
		return
	}

	// 2 <--- 1
	if box1InCircuit && !box2InCircuit {
		data.CircuitsTable[minPair.Box2] = circuit1
		// fmt.Printf("\tBox 2 is not in a circuit, joining id: %d\n", circuit1)
		return
	}

	// 1 <--- 2
	if box2InCircuit && !box1InCircuit {
		data.CircuitsTable[minPair.Box1] = circuit2
		// fmt.Printf("\tBox 1 is not in a circuit, joining id: %d\n", circuit2)
		return
	}

	// Both in a circuit
	if box1InCircuit && box2InCircuit {
		if circuit1 != circuit2 {
			// fmt.Printf("\tMerge circuits %d and %d\n", circuit1, circuit2)
			// Merge circuit 1 and 2
			// Generate a new circuit newCircuitId
			newCircuitId := data.NextCircuitId
			data.NextCircuitId++
			// Update any box in circuit 1 or 2
			for key, val := range data.CircuitsTable {
				if val == circuit1 || val == circuit2 {
					data.CircuitsTable[key] = newCircuitId
				}
			}
		} else {
			// fmt.Printf("\tBoth in the same circuit --> do nothing")
		}
	}
}

func IsUniqueCircuit(data *BoxesData) bool {
	last_id := -1

	for _, id := range data.CircuitsTable {
		if last_id == -1 {
			last_id = id
		}

		if id != last_id {
			return false
		}
	}

	return true
}

func Top3Circuits(data BoxesData) []uint64 {
	maxQueue := make(MaxQueue, 0)
	heap.Init(&maxQueue)

	circuits := make(map[int]int)

	for _, circuitId := range data.CircuitsTable {

		if _, found := circuits[circuitId]; !found {
			circuits[circuitId] = 0
		}
		circuits[circuitId]++
	}

	for circuitId, quantity := range circuits {
		heap.Push(&maxQueue, &Circuit{id: circuitId, quantity: quantity})
	}

	result := make([]uint64, 0, 3)
	for range 3 {
		circuit := heap.Pop(&maxQueue).(*Circuit)
		fmt.Printf("Circuit %d: %d\n", circuit.id, circuit.quantity)
		result = append(result, uint64(circuit.quantity))
	}

	return result
}

func Part1(data []IVec3) uint64 {

	// Buscar el par de cajas mas cercanos
	// Mergearlos en un circuito

	/*
		Associates each box coordinates with a circuit id
		If box is not in the map or its id is 0, the box is not in a circuit
	*/

	d := BoxesData{
		Boxes:         data,
		CircuitsTable: make(map[IVec3]int),
		DistanceCache: make(map[Pair]float64),
		NextCircuitId: 1,
	}

	pairs := GetSortedPairs(&d)
	var minPair Pair

	for range 1000 {
		minPair = pairs[0]
		pairs = pairs[1:]

		if IsPairLinked(&d, minPair) {
			// the pair already is linked in the same circuit
			continue
		}

		MergeCircuits(minPair, &d)
	}

	total := uint64(1)
	for _, c := range Top3Circuits(d) {
		total *= c
	}
	return total
}

func Part2(data []IVec3) int {

	d := BoxesData{
		Boxes:         data,
		CircuitsTable: make(map[IVec3]int),
		DistanceCache: make(map[Pair]float64),
		NextCircuitId: 1,
	}

	pairs := GetSortedPairs(&d)
	var minPair Pair
	total := 1

	for len(pairs) > 0 {
		minPair = pairs[0]
		pairs = pairs[1:]

		if IsPairLinked(&d, minPair) {
			// the pair already is linked in the same circuit
			continue
		}

		MergeCircuits(minPair, &d)
		if IsUniqueCircuit(&d) {
			total = minPair.Box1.X * minPair.Box2.X
		}
	}

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
	var part2 int = Part2(data)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
