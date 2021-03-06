package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type memory []int
type channel chan int

func main() {
	items := []string{
		"jam",
		"bowl of rice",
		"antenna",
		"manifold",
		"hypercube",
		"dehydrated water",
		"candy cane",
		"dark matter",
	}

	// 2147502592
	//Taking jam
	//Taking antenna
	//Taking hypercube
	//Taking dehydrated water
	//Taking candy cane

	memory, in, out, stop := load()
	go func() {
		reader := bufio.NewReader(os.Stdin)
		commands := []string{}

		// For brute-forcing weight-lock
		bruteForce := false
		bfPosition := 0
		combinations := make([][]string, 256)

		for {
			// Display drone messages.
			messageShown := false
			for !messageShown {
				for len(out) > 0 {
					ch := <-out
					fmt.Print(string(ch))
					messageShown = true
				}
			}

			// Wait for input
			var input string
			if !bruteForce {
				bs, _, _ := reader.ReadLine()
				input = string(bs)
			} else {
				fmt.Println(strings.Repeat("-", 80))
				// Check that particular BF Position.
				// Drop everything.
				for _, value := range items {
					in.send("drop " + value)
				}
				// Take only the ones in the combination.
				for _, item := range combinations[bfPosition] {
					fmt.Println("Taking", item)
					in.send("take " + item)
				}
				bfPosition++

				// Walk west.
				in.send("west")

				// Examine output
				wait()
			}

			// Own commands
			if input == "brute-force" {
				bruteForce = true
				// Create a list of all combinations 2^8 = 256.
				for i := 0; i < 256; i++ {
					combinations[i] = []string{}
					// Check if the specific item is in this combination
					for pos, item := range items {
						isSet := i&(1<<pos) > 0
						if isSet {
							combinations[i] = append(combinations[i], item)
						}
					}
				}
				fmt.Println("breakpoint")
			}
			if input == "save" {
				file, _ := os.Create("savegame")
				for _, value := range commands {
					_, _ = file.WriteString(value + "\n")
				}
				_ = file.Close()
				fmt.Println("SAVED")
			}
			if input == "load" {
				// NOTE: Will only work at the beginning
				file, _ := os.Open("savegame")
				bs, _ := ioutil.ReadAll(file)
				inputs := strings.Split(string(bs), "\n") // Use scanner instead?
				for _, command := range inputs {
					if command != "save" {
						in.send(command)
					}
				}
				commands = inputs
				file.Close()
				fmt.Println("LOADED")
			}

			parts := strings.Split(input, " ")
			switch parts[0] {
			case "n":
				input = "north"
			case "s":
				input = "south"
			case "w":
				input = "west"
			case "e":
				input = "east"
			case "t":
				input = "take " + strings.Join(parts[1:], " ")
			case "d":
				input = "drop " + strings.Join(parts[1:], " ")
			}
			if len(input) > 0 {
				in.send(input)
				commands = append(commands, input)
			}
		}
	}()
	compute(memory, in, out, stop)
}

func compute(memory memory, in channel, out channel, stop *bool) {
	relBase := 0

	for ip := 0; ip < len(memory); {
		mem := memory[ip]
		opcode := mem % 100
		r1 := mem / 100 % 10
		r2 := mem / 1000 % 10
		r3 := memory[ip] / 10000 % 10

		switch opcode {
		case 1:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if r2 == 2 {
				m2 = memory[memory[ip+2]+relBase]
			}
			var m3 int
			if r3 == 0 {
				m3 = memory[ip+3]
			}
			if r3 == 2 {
				m3 = memory[ip+3] + relBase
			}
			memory[m3] = m1 + m2
			ip += 4
		case 2:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if r2 == 2 {
				m2 = memory[memory[ip+2]+relBase]
			}
			var m3 int
			if r3 == 0 {
				m3 = memory[ip+3]
			}
			if r3 == 1 {
				panic("TODO")
			}
			if r3 == 2 {
				m3 = memory[ip+3] + relBase
			}
			memory[m3] = m1 * m2
			ip += 4
		case 3:
			var num int
			num = <-in
			if r1 == 0 {
				memory[memory[ip+1]] = num
			}
			if r1 == 1 {
				panic("How should this work?")
			}
			if r1 == 2 {
				memory[memory[ip+1]+relBase] = num
			}
			ip += 2
		case 4:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			out <- m1
			ip += 2
		case 5:
			// Jump if true.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if r2 == 2 {
				m2 = memory[memory[ip+2]+relBase]
			}
			if m1 != 0 {
				ip = m2
			} else {
				ip += 3
			}
		case 6:
			// Jump if false.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if r2 == 2 {
				m2 = memory[memory[ip+2]+relBase]
			}
			if m1 == 0 {
				ip = m2
			} else {
				ip += 3
			}
		case 7:
			// Less than.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if r2 == 2 {
				m2 = memory[memory[ip+2]+relBase]
			}
			var m3 int
			if r3 == 0 {
				m3 = memory[ip+3]
			}
			if r3 == 1 {
				panic("TODO")
			}
			if r3 == 2 {
				m3 = memory[ip+3] + relBase
			}
			if m1 < m2 {
				memory[m3] = 1
			} else {
				memory[m3] = 0
			}
			ip += 4
		case 8:
			// Equals.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if r2 == 2 {
				m2 = memory[memory[ip+2]+relBase]
			}
			var m3 int
			if r3 == 0 {
				m3 = memory[ip+3]
			}
			if r3 == 1 {
				panic("TODO")
			}
			if r3 == 2 {
				m3 = memory[ip+3] + relBase
			}
			if m1 == m2 {
				memory[m3] = 1
			} else {
				memory[m3] = 0
			}
			ip += 4
		case 9:
			// Relative base adjustment.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			if r1 == 2 {
				m1 = memory[memory[ip+1]+relBase]
			}
			relBase += m1
			ip += 2
		case 99:
			for len(out) > 0 {
				time.Sleep(time.Millisecond * 100)
			}
			*stop = true
			return
		default:
			panic("Unknown opcode:" + strconv.Itoa(opcode))
		}
	}
}

func load() (memory, channel, channel, *bool) {
	const MemorySize = 1000000
	const ChannelSize = 16384

	var memory []int
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), ",")
	memory = make([]int, MemorySize)
	for idx, val := range lines {
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		memory[idx] = i
	}
	in := make(chan int, ChannelSize)
	out := make(chan int, ChannelSize)
	stop := false
	return memory, in, out, &stop
}

func (c channel) send(msg string) {
	trimmed := strings.Trim(msg, " ")
	for i := 0; i < len(trimmed); i++ {
		r := msg[i]
		c <- int(r)
	}
	c <- 10 // Newline
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}
