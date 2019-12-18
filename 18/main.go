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
	fmt.Println(keys)

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
	fmt.Println(doors)
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
