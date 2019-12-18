package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	input := load()

	// Find starting point.
	var x, y int
	withInput(input, func(_x, _y, value int) {
		if value == '@' {
			x = _x
			y = _y
		}
	})

	fmt.Println(x, y)
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
		colMem := make([]int, len(row))
		input = append(input, colMem)
		for x, col := range row {
			input[y][x] = int(col)
		}
	}

	return input
}
