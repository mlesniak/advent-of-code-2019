package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	// Part 1.
	memory := load()
	compute(memory)
	showResult(memory)
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
	lines := strings.Split(string(bytes), ",")
	var memory []int
	for _, val := range lines {
		i, _ := strconv.Atoi(val)
		memory = append(memory, i)
	}
	return memory
}
