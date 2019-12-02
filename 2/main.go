package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
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
