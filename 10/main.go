package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	memory := load()
	fmt.Println(memory)
}

type coordinate struct {
	x int
	y int
}

func load() map[coordinate]bool {
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	asteroids := make(map[coordinate]bool)
	for row, line := range lines {
		fmt.Println(row, line)
		for col, char := range line {
			if char == '#' {
				asteroids[coordinate{x: col, y: row}] = true
			}
		}
	}
	return asteroids
}
