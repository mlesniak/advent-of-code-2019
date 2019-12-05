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
	for i := 0; i < len(memory); {
		switch memory[i] {
		case 1:
			memory[memory[i+3]] = memory[memory[i+1]] + memory[memory[i+2]]
			i += 4
		case 2:
			memory[memory[i+3]] = memory[memory[i+1]] * memory[memory[i+2]]
			i += 4
		case 3:
			var num int
			fmt.Print("? ")
			_, err := fmt.Scanf("%d", &num)
			if err != nil {
				panic(err)
			}
			memory[memory[i+1]] = num
			i += 2
		case 4:
			fmt.Println(memory[memory[i+1]])
			i += 2
		case 99:
			return
		}
	}
}

func load() []int {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), ",")
	fmt.Println(lines)
	var memory []int
	for _, val := range lines {
		i, _ := strconv.Atoi(val)
		memory = append(memory, i)
	}
	return memory
}
