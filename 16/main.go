package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	//x := []int{9,8,7,6,5,4,3,2,1,0,9,8,7,6,5,4,3,2,1,0}
	//fmt.Println(x[7:7+8])
	//return

	times := 10000
	base := load()
	input := make([]int, 0) // Initialize before?
	fmt.Println(times * len(base))

	// Create list repeated 10000 times.
	for i := 0; i < times; i++ {
		input = append(input, base...)
	}

	// [5 9 7 5 0 5 3]
	//offset := input[:7]
	//fmt.Println(offset)
	offset := 5975053

	steps := 100
	for i := 0; i < steps; i++ {
		log.Println("Tick", i)
		output := compute(offset, input)
		input = output
	}

	fmt.Println(input[offset : offset+8])
}

func compute(offset int, input []int) []int {
	output := make([]int, len(input))

	sum := 0
	for pos := len(output) - 1; pos >= offset; pos-- {
		sum += input[pos]
		output[pos] = sum
		output[pos] %= 10
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
