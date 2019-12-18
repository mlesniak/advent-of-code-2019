package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

func main() {
	view := load()

	// Find starting point.
	var x, y int
	withInput(view, func(_x, _y, value int) {
		if value == '@' {
			x = _x
			y = _y
		}
	})

	// Idea: explore whole lab with backtracking until we have found all keys; Dijkstra would be another option?
	// Ignore doors for now. Directions:
	//
	//          0
	//        2   3
	//          1
	foundKeys := make(map[int]bool)
	path := []direction{}
	backtrack(view, foundKeys, path, x, y)
}

type direction int

const (
	up = iota
	down
	left
	right
)

const minDirection = up
const maxDirection = right

func (d direction) String() string {
	switch d {
	case up:
		return "U"
	case down:
		return "D"
	case left:
		return "L"
	case right:
		return "R"
	}

	panic("Unknown direction")
}

func backtrack(view [][]int, foundKeys map[int]bool, path []direction, x int, y int) {
	if len(path) > 10 {
		return
	}

	//fmt.Println("\n" + strings.Repeat("-", 40))
	//fmt.Println("x=", x, ", y=", y, ", value=", string(view[y][x]))
	//fmt.Println("path", path)
	//fmt.Println("foundKeys", foundKeys)
	//wait()

	dxs := []int{-1, 1, 0}
	dys := []int{-1, 1, 0}

	// Check if we stand on a key.
	if view[y][x] >= 'a' && view[y][x] <= 'z' {
		fmt.Println("Found key", view[y][x])
		foundKeys[view[y][x]] = true
		view[y][x] = '.'
	}
	if len(foundKeys) == 2 {
		fmt.Println("FOUND")
		fmt.Println(path)
		return
	}

	for _, dx := range dxs {
		for _, dy := range dys {
			// Ignore diagnoal paths.
			abs := math.Abs(float64(dx)) + math.Abs(float64(dy))
			if abs > 1 || abs == 0 {
				continue
			}

			// Since we have a border around the whole maze, we do not need to check for negative or too large indices.
			ny := y + dy
			nx := x + dx
			val := view[ny][nx]

			onStart := false
			if view[ny][nx] == '@' {
				onStart = true
			}

			onDoor := false
			if view[ny][nx] >= 'A' && view[ny][nx] <= 'Z' {
				onDoor = true
			}

			onKey := false
			if view[ny][nx] >= 'a' && view[ny][nx] <= 'z' {
				onKey = true
			}

			// Check if this is not the opposite direction.
			direction := fromDelta(dx, dy)
			//opposite := len(path) > 0 && direction == oppositeDirection(path[len(path)-1])
			//opposite := false

			// Check if path is free.
			freeField := val == '.'

			if freeField || onDoor || onKey || onStart {
				// Add to path and backtrack.
				nPath := append(path, direction)
				copyKeys := make(map[int]bool)
				for key, value := range foundKeys {
					copyKeys[key] = value
				}
				backtrack(view, copyKeys, nPath, nx, ny)
				//fmt.Println("Backtrack from", nPath, "to", path)
			}
		}
	}
}

func oppositeDirection(d direction) direction {
	switch d {
	case up:
		return down
	case down:
		return up
	case left:
		return right
	case right:
		return left
	}

	panic("Unknown direction")
}

func fromDelta(dx int, dy int) direction {
	switch {
	case dx == 0 && dy == 1:
		return down
	case dx == 0 && dy == -1:
		return up
	case dx == 1 && dy == 0:
		return right
	case dx == -1 && dy == 0:
		return left
	}

	panic("Unknown dx/dy")
}

func withInput(input [][]int, f func(x, y, value int)) {
	for row := range input {
		for col := range input[row] {
			f(col, row, input[row][col])
		}
	}
}

func load() [][]int {
	input := make([][]int, 0)

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for y, row := range lines {
		fmt.Println(row)
		colMem := make([]int, len(row))
		input = append(input, colMem)
		for x, col := range row {
			input[y][x] = int(col)
		}
	}

	return input
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}
