package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type channel chan int

func main() {
	width := 1000
	height := 100
	buffer := [][]int{}

	// Initially, fill buffer with height lines.
	for y := 0; y < height; y++ {
		fmt.Println("Init", y)
		buffer = append(buffer, readLine(width, y))
	}
	//show(buffer, height)

	// Note that we will not find rectangles in the first 50 lines with our approach. ¯\_(ツ)_/¯
	const rectangleSize = 100
	y := height
	fmt.Println(y)
	for {
		// Check the last-height line for a matching width.
		top := height - rectangleSize
		line := buffer[top]
		// Find a sequence of 1s with given length.
		for i := 0; i < len(line); i++ {
			if line[i] == 1 {
				if i+rectangleSize < len(line) && line[i+rectangleSize] == 1 {
					// Check bottom
					if buffer[len(buffer)-1][i] == 1 {
						// Found.
						//show(buffer, y)
						fmt.Println(i, top)
						fmt.Print(top, " ")
						//showLine(line)
						solution := i*10000 + top
						fmt.Println(solution)
						os.Exit(0)
					}
				}
			}
		}

		// Update buffer, remove first line, update last one.
		y++
		fmt.Println(y)
		newLine := readLine(width, y)
		buffer = append(buffer[1:], newLine)
	}
}

//.......................######.....................
//12345678901234567890123456

func showLine(line []int) {
	for _, value := range line {
		var c string
		switch value {
		case 0:
			c = "."
		case 1:
			c = "#"
		}
		fmt.Print(c)
	}
	fmt.Println()
}

func readLine(width int, y int) []int {
	buffer := make([]int, width)

	var w sync.Mutex
	for x := 0; x < width; x++ {
		w.Lock()
		go func(x, y int) {
			memory, in, out, stop := load()
			go func() {
				in <- x
				in <- y
				c := <-out
				buffer[x] = c
			}()
			compute(memory, in, out, stop)
			for !*stop {
				// Wait...
			}
			w.Unlock()
		}(x, y)
	}
	w.Lock()

	return buffer
}

func add() {
	//finished := make(chan bool, size*size)
	//for y := 0; y < size; y++ {
	//	for x := 0; x < size; x++ {
	//		go func(x, y int) {
	//			memory, in, out, stop := load()
	//			go func() {
	//				in <- x
	//				in <- y
	//				c := <-out
	//				if y < size && x < size {
	//					// Strange that I have to check this...
	//					buffer[y][x] = c
	//				}
	//			}()
	//			compute(memory, in, out, stop)
	//			for !*stop {
	//				// Wait...
	//			}
	//			finished <- true
	//		}(x, y)
	//	}
	//}
	//
	//for len(finished) < size*size {
	//	time.Sleep(time.Millisecond * 100)
	//}
}

func show(view [][]int, maxY int) {
	delta := maxY - len(view) + 1
	for row := range view {
		fmt.Print(row+delta, " ")
		for col := range view[row] {
			var c string
			switch view[row][col] {
			case 0:
				c = "."
			case 1:
				c = "#"
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func compute(memory []int, in channel, out channel, stop *bool) {
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

func load() ([]int, channel, channel, *bool) {
	const MemorySize = 1000000
	const ChannelSize = 128

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), ",")
	memory := make([]int, MemorySize)
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
