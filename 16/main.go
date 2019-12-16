package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	times := 10000
	base := load()
	input := make([]int, 0) // Initialize before?
	fmt.Println(times * len(base))

	// Create list repeated 10000 times.
	for i := 0; i < times; i++ {
		input = append(input, base...)
	}

	steps := 100
	for i := 0; i < steps; i++ {
		log.Println("Tick", i)
		output := compute(input)
		input = output
	}

	fmt.Println(input[:8])
}

func compute(input []int) []int {
	output := make([]int, len(input))
	for pos := 0; pos < len(output); pos++ {
		if pos%10 == 0 {
			log.Println("Pos=", pos)
		}
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
