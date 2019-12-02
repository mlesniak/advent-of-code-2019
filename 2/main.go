package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func main() {
	// Part 1.
	//memory := load()
	//compute(memory)
	//showResult(memory)

	// Part 2.
	result := math.MaxInt64
	const goal = 19690720
loop:
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			fmt.Println(noun, verb)
			memory := load()
			memory[1] = noun
			memory[2] = verb
			compute(memory)
			if memory[0] == goal {
				result = 100*noun + verb
				break loop
			}
		}
	}
	if result == math.MaxInt64 {
		fmt.Println("Probably no solution found, code buggy?")
	} else {
		fmt.Println(result)
	}
}

func showResult(memory []int) {
	fmt.Println(memory[0])
}

func compute(memory []int) {
	for i := 0; i < len(memory); i += 4 {
		switch memory[i] {
		case 1:
			memory[memory[i+3]] = memory[memory[i+1]] + memory[memory[i+2]]
		case 2:
			memory[memory[i+3]] = memory[memory[i+1]] * memory[memory[i+2]]
		case 99:
			break
		}
	}
}

func load() []int {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	// For testing, see https://adventofcode.com/2019/day/2.
	//bytes := []byte("1,9,10,3,2,3,11,0,99,30,40,50")
	lines := strings.Split(string(bytes), ",")
	var memory []int
	for _, val := range lines {
		i, _ := strconv.Atoi(val)
		memory = append(memory, i)
	}
	return memory
}
