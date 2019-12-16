package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	//input := load()
	//fmt.Println(input)

	// Repeat by modulo operation
	pattern := computePattern(2)
	fmt.Println(pattern)
}

func computePattern(position int) []int {
	basePattern := []int{0, 1, 0, -1}

	var pattern []int
	for _, digit := range basePattern {
		for i := 0; i <= position; i++ {
			pattern = append(pattern, digit)
		}
	}

	return pattern
}

func load() []int {
	var input []int

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		for _, digit := range line {
			i, err := strconv.Atoi(string(digit))
			if err != nil {
				panic(err)
			}
			input = append(input, i)
		}
	}

	return input
}
