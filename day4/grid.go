package main

import (
	"fmt"
	"strings"
)

type Grid[T any] struct {
	Data    []T
	Rows    int
	Columns int
}

type Position struct {
	Row    int
	Column int
}

func (p Position) Add(other Position) Position {
	return Position{Row: p.Row + other.Row, Column: p.Column + other.Column}
}

var up = Position{Row: -1, Column: 0}
var down = Position{Row: +1, Column: 0}
var left = Position{Row: 0, Column: -1}
var right = Position{Row: 0, Column: +1}

var upLeft = Position{Row: -1, Column: -1}
var downLeft = Position{Row: +1, Column: -1}
var upRight = Position{Row: -1, Column: +1}
var downRight = Position{Row: +1, Column: +1}

func NewGrid[T any](rows int, columns int) *Grid[T] {
	g := new(Grid[T])

	size := rows * columns

	g.Rows = rows
	g.Columns = columns
	g.Data = make([]T, size)

	return g
}

func (g *Grid[T]) Get4Neighbours(pos Position) []Position {
	directions := []Position{up, down, left, right}
	valid := []Position{}

	for _, dir := range directions {
		neigh := pos.Add(dir)

		if g.IsValidPosition(neigh) {
			valid = append(valid, neigh)
		}
	}

	return valid
}

func (g *Grid[T]) Get8Neighbours(pos Position) []Position {
	directions := []Position{up, down, left, right, upLeft, upRight, downLeft, downRight}
	valid := []Position{}

	for _, dir := range directions {
		neigh := pos.Add(dir)

		if g.IsValidPosition(neigh) {
			valid = append(valid, neigh)
		}
	}

	return valid
}

func (g *Grid[T]) String() string {
	var strBuilder strings.Builder

	for r := range g.Rows {
		for c := range g.Columns {
			ptr := g.GetPtrAt(Position{Row: r, Column: c})
			strBuilder.WriteString(fmt.Sprint(*ptr))
		}
		strBuilder.WriteRune('\n')
	}

	return strBuilder.String()

}

func (g *Grid[T]) Index(pos Position) int {
	return pos.Row*g.Columns + pos.Column
}

func (g *Grid[T]) IsValidPosition(pos Position) bool {
	return pos.Row >= 0 && pos.Column >= 0 &&
		pos.Row < g.Rows && pos.Column < g.Columns
}

func (g *Grid[T]) SetAt(pos Position, value T) {
	if !g.IsValidPosition(pos) {
		return
	}
	g.Data[g.Index(pos)] = value
}

func (g *Grid[T]) GetPtrAt(pos Position) *T {
	if !g.IsValidPosition(pos) {
		return nil
	}

	return &g.Data[g.Index(pos)]
}
