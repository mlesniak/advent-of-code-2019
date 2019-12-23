package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type channel chan int

func main() {
	// Use binary search to find a line.

	rectangleSize := 100
	fmt.Println("Rectangle size", rectangleSize)
	width := 1000
	min := 0
	max := 2000

	pos := 0
	oldPos := 1
	for pos != oldPos {
		oldPos = pos
		pos = (min + max) / 2
		log.Println("Examining pos=", pos)
		top := readLine(width, pos)

		// Check if this top has enough ones.
		found := false
		for x, value := range top {
			if value == 1 && x < len(top)-1-rectangleSize {
				if top[x+rectangleSize] == 1 {
					// Looks promising. Could be a top corner, hence look at the bottom left corner.
					bottom := readLine(width, pos+rectangleSize)
					if bottom[x] == 1 {
						// Corner found. Candidate.
						solution := x*10000 + pos
						fmt.Println("SOLUTION", x, pos, solution)

						//w := x + rectangleSize
						//for y := pos; y < pos+rectangleSize; y++ {
						//		l := readLine(w, y)
						//		showLine(x-5, l)
						//}

						found = true
						break
					}
				}
			}
		}

		count := 0
		for _, value := range top {
			if value == 1 {
				count++
			}
		}
		fmt.Println("pos=", pos, "; count(1)=", count)

		if count < rectangleSize {
			min = pos
		} else if found {
			max = pos
		} else {
			min = pos
		}
	}
}

func showLine(startAt int, line []int) {
	for idx, value := range line {
		if idx <= startAt {
			continue
		}
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
