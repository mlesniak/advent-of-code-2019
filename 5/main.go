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
	//fmt.Println(memory)
	fmt.Println(memory[0])
}

func compute(memory []int) {
	for i := 0; i < len(memory); {
		opcode := memory[i] % 100
		r1 := memory[i] / 100 % 10
		r2 := memory[i] / 1000 % 10
		//r3 := memory[i] / 10000 % 10

		switch opcode {
		case 1:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[i+1]]
			}
			if r1 == 1 {
				m1 = memory[i+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[i+2]]
			}
			if r2 == 1 {
				m2 = memory[i+2]
			}
			var m3 int
			m3 = memory[memory[i+3]]
			memory[m3] = m1 + m2
			i += 4
		case 2:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[i+1]]
			}
			if r1 == 1 {
				m1 = memory[i+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[i+2]]
			}
			if r2 == 1 {
				m2 = memory[i+2]
			}
			var m3 int
			m3 = memory[i+3]
			memory[m3] = m1 * m2
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
		default:
			panic("Unknown opcode:" + strconv.Itoa(opcode))
		}
	}
}

func load() []int {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), ",")
	//fmt.Println(lines)
	var memory []int
	for _, val := range lines {
		i, _ := strconv.Atoi(val)
		memory = append(memory, i)
	}
	return memory
}
