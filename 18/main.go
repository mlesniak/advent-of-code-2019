package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

	// Find coordinates of every key.
	keys := make(map[int]coordinate)
	for c := 'a'; c <= 'z'; c++ {
		withInput(view, func(x, y, value int) {
			ic := int(c)
			if value == ic {
				keys[ic] = coordinate{x, y}
			}
		})
	}

	// Find coordinates of every door, necessary for opening (=removing) a door.
	doors := make(map[int]coordinate)
	for c := 'A'; c <= 'Z'; c++ {
		withInput(view, func(x, y, value int) {
			ic := int(c)
			if value == ic {
				doors[ic] = coordinate{x, y}
			}
		})
	}

	// Implement simple BFS to find all possible keys.
	bfs(view, x, y, 'a')
}

type path struct {
	position coordinate
	length   int
}

func bfs(view [][]int, x int, y int, key int) int {
	start := coordinate{x, y}
	candidates := []path{{start, 0}}
	history := make(map[coordinate]bool)

	for len(candidates) > 0 {
		//wait()
		//fmt.Println(candidates)
		candidate := candidates[0]
		position := candidate.position
		//fmt.Println("- Looking at", candidate)
		candidates = candidates[1:]

		// Ignore already visited.
		_, visited := history[position]
		if visited {
			continue
		}
		history[position] = true

		// Ignore out-of-fields.
		if position.x < 0 || position.x > len(view[0]) || position.y < 0 || position.y > len(view) {
			continue
		}
		// Ignore walls.
		if view[position.y][position.x] == '#' {
			continue
		}
		// Ignore doors.
		if view[position.y][position.x] >= 'A' && view[position.y][position.x] <= 'Z' {
			continue
		}

		// Check if found.
		if view[position.y][position.x] == key {
			fmt.Println("Found", candidate)
			return candidate.length
		}

		// Generate new candidates.
		addCandidate(history, &candidates, position.x+1, position.y, candidate.length)
		addCandidate(history, &candidates, position.x-1, position.y, candidate.length)
		addCandidate(history, &candidates, position.x, position.y+1, candidate.length)
		addCandidate(history, &candidates, position.x, position.y-1, candidate.length)
	}

	return -1
}

func addCandidate(history map[coordinate]bool, candidates *[]path, x int, y int, len int) {
	if history[coordinate{x, y}] == false {
		p := path{coordinate{x, y}, len + 1}
		*candidates = append(*candidates, p)
	}
}

type coordinate struct {
	x, y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("%d/%d", c.x, c.y)
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
