package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	input := load()
	//fmt.Println(input)

	//input := []int{1, 2, 3, 4, 5, 6, 7, 8}
	//fmt.Println(input)

	steps := 100
	for i := 0; i < steps; i++ {
		output := compute(input)
		fmt.Println(output)
		input = output
	}

	fmt.Println(input[:8])
}

func compute(input []int) []int {
	output := make([]int, len(input))
	for pos := 0; pos < len(output); pos++ {
		pattern := computePattern(pos)
		for i := 0; i < len(input); i++ {
			factor := getFactor(pattern, i)
			output[pos] += input[i] * factor
		}
		output[pos] %= 10
		if output[pos] < 0 {
			output[pos] *= -1
		}
	}
	return output
}

//1, 0, -1, 0, 1, 0, -1, 0, 1, 0, -1,
func getFactor(pattern []int, pos int) int {
	if pos+1 < len(pattern) {
		return pattern[pos+1]
	}

	return pattern[(pos+1)%len(pattern)]
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
