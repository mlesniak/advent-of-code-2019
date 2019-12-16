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

	input := []int{1, 2, 3, 4, 5, 6, 7, 8}
	output := make([]int, len(input))
	for i := 0; i < len(output); i++ {
		fmt.Println("Computing for position:", i)
		output[i] = input[i]
	}
	fmt.Println(output)
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
