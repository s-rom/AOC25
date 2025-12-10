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

func IsEdgeHorizontal(a IVec2, b IVec2) bool {
	return a.Y == b.Y
}

func IsEdgeVertical(a IVec2, b IVec2) bool {
	return a.X == b.X
}

func IsPointInEdge(p IVec2, e1 IVec2, e2 IVec2) bool {
	if IsEdgeHorizontal(e1, e2) {
		minX := min(e1.X, e2.X)
		maxX := max(e1.X, e2.X)
		return p.Y == e1.Y && p.X <= maxX && p.X >= minX
	} else {
		minY := min(e1.Y, e2.Y)
		maxY := max(e1.Y, e2.Y)
		return p.X == e1.X && p.Y >= minY && p.Y <= maxY
	}
}

func TestIntersectionVH(v1 IVec2, v2 IVec2, h1 IVec2, h2 IVec2) bool {
	vX := v1.X // Vertical edge: fixed x
	hY := h1.Y // Horizontal edge: fixed y

	hMaxX := max(h1.X, h2.X)
	hMinX := min(h1.X, h2.X)

	vMinY := min(v1.Y, v2.Y)
	vMaxY := max(v1.Y, v2.Y)

	if vX < hMaxX && vX > hMinX && hY > vMinY && hY < vMaxY {
		return true
	}

	return false
}

func TestEdgeIntersectsPolygon(e1 IVec2, e2 IVec2, vertices []IVec2) bool {

	for idx := range vertices {
		polyA := vertices[idx]
		polyB := vertices[(idx+1)%len(vertices)]

		if IsEdgeHorizontal(polyA, polyB) && IsEdgeVertical(e1, e2) {
			if TestIntersectionVH(e1, e2, polyA, polyB) {
				return true
			}
		} else if IsEdgeVertical(polyA, polyB) && IsEdgeHorizontal(e1, e2) {
			if TestIntersectionVH(polyA, polyB, e1, e2) {
				return true
			}
		}
	}

	return false
}

func IsPointInAnyEdge(p IVec2, vertices []IVec2) bool {

	for idx := range len(vertices) {

		a := vertices[idx]
		b := vertices[(idx+1)%len(vertices)]

		if IsPointInEdge(p, a, b) {
			return true
		}
	}

	return false

}

func IsPointInsidePoligon(p IVec2, vertices []IVec2) bool {

	intersections := 0
	for idx := range len(vertices) {

		a := vertices[idx]
		b := vertices[(idx+1)%len(vertices)]

		if IsEdgeVertical(a, b) && a.X > p.X {
			minY := min(a.Y, b.Y)
			maxY := max(a.Y, b.Y)
			if p.Y >= minY && p.Y < maxY {
				intersections++
			}
		}

	}

	inside := (intersections%2 == 1) // inside if odd
	return inside
}
func Part2(data []IVec2) uint64 {
	total := uint64(0)
	pointInsideCache := make(map[IVec2]bool)

	for idxA, a := range data {
		for idxB := idxA + 1; idxB < len(data); idxB++ {
			b := data[idxB]
			minX := min(a.X, b.X)
			maxX := max(a.X, b.X)
			minY := min(a.Y, b.Y)
			maxY := max(a.Y, b.Y)

			// Corners of the rectangle in clockwise order
			corners := []IVec2{
				{X: minX, Y: minY},
				{X: maxX, Y: minY},
				{X: maxX, Y: maxY},
				{X: minX, Y: maxY},
			}

			// Check if edge a-b intersects any of the edges
			//	If intersection => invalid rectangle
			validRectangle := true
			for cornerIdx := range corners {
				a := corners[cornerIdx]
				b := corners[(cornerIdx+1)%len(corners)]
				if TestEdgeIntersectsPolygon(a, b, data) {
					validRectangle = false
					break
				}
			}

			if !validRectangle {
				continue
			}

			rectangleInside := true
			for _, corner := range corners {
				var cornerInside bool
				val, found := pointInsideCache[corner]

				if found {
					cornerInside = val
				} else {
					inEdge := IsPointInAnyEdge(corner, data)
					cornerInside = inEdge || (!inEdge && IsPointInsidePoligon(corner, data))
					pointInsideCache[corner] = cornerInside
				}

				if !cornerInside {
					rectangleInside = false
					break
				}
			}

			if rectangleInside {
				area := GetRectArea(a, b)
				if area > total {
					total = area
				}
			}
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
	var part2 uint64 = Part2(data)

	fmt.Printf("Part 1: %d\nPart 2: %d\n", part1, part2)

}
