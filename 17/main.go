package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type channel chan int

func (c channel) send(msg string) {
	if len(msg) > 20 {
		panic(fmt.Sprintf("Message too long len(%s)=%d", msg, len(msg)))
	}
	for i := 0; i < len(msg); i++ {
		r := msg[i]
		c <- int(r)
	}
	c <- 10 // Newline
}

func main() {
	memory, in, out, stop := load()

	width := 50
	height := 50
	view := make([][]int, height)
	for row := range view {
		view[row] = make([]int, width)
	}
	go func() {
		renderAndStoreView(stop, in, out, view)

		path := computePath(view)
		fmt.Println(path)

		//in.send("A,B,C")       // Code
		//in.send("R,8,L,10,R,8,R,12") // A
		//in.send("R,8,L,8,L,12,R,8")       // B
		//in.send("L,10,R,8,L,12,L,10")       // C
		//in.send("n")
		//for {
		//	if *stop {
		//		break
		//	}
		//	renderAndStoreView(stop, in, out, view)
		//
		//	// Final result
		//	//num := <- out
		//	//if num > 255 {
		//	//	fmt.Println(num)
		//	//}
		//}
	}()
	compute(memory, in, out, stop)
}

func computePath(view [][]int) string {
	sx, sy, _ := findView(view, func(x int, y int, value int) bool {
		return value == int('^')
	})

	// We cheat a bit by setting first value.
	path := "R,"
	dx := 1
	dy := 0
	x := sx
	y := sy
	for {
		//fmt.Println(path)
		//fmt.Println(x+dx, y+dy)

		steps := 0
		// Check if the next step in a direction is still a tile?
		for y+dy < len(view) && y+dy >= 0 && view[y+dy][x+dx] == '#' {
			steps++
			x += dx
			y += dy
		}
		path += strconv.Itoa(steps) + ","

		// Try turning left and check view.
		dx, dy = left(dx, dy)
		if y+dy < len(view) && y+dy >= 0 && view[y+dy][x+dx] == '#' {
			path += "L,"
		} else {
			dx, dy = left(dx, dy)
			dx, dy = left(dx, dy)
			if view[y+dy][x+dx] == '#' {
				path += "R,"
			} else {
				// Finished
				break
			}
		}
	}

	return path
}

func right(dx int, dy int) (int, int) {
	d1, d2 := left(dx, dy)
	d1, d2 = left(d1, d2)
	return left(d1, d2)
}

func left(dx int, dy int) (int, int) {
	switch {
	case dx == 1 && dy == 0:
		return 0, -1
	case dx == 0 && dy == -1:
		return -1, 0
	case dx == -1 && dy == 0:
		return 0, 1
	case dx == 0 && dy == 1:
		return 1, 0
	}

	panic("Should not happen")
}

func findView(view [][]int, f func(int, int, int) bool) (int, int, int) {
	for y := range view {
		for x := range view[y] {
			abort := f(x, y, view[y][x])
			if abort {
				return x, y, view[y][x]
			}
		}
	}

	return -1, -1, -1
}

func renderAndStoreView(stop *bool, in chan int, out chan int, view [][]int) {
	y := 0
	x := 0
	for y <= 37 {
		c := <-out
		if c > 255 {
			fmt.Println(c)
		} else {
			fmt.Print(string(c))
		}

		if c == 10 {
			y++
			x = 0
		} else {
			view[y][x] = c
			x++
		}
	}
}

func findIntersections(view [][]int, height int, width int) {
	// Find all intersections
	sum := 0
	for y := range view {
		// In our case, there are no 'T'-shapes at the edges, simplifying computation.
		if y == 0 || y == height-1 {
			continue
		}
		for x := range view[y] {
			if x == 0 || x == width-1 {
				continue
			}

			// Find only #
			scaffold := int('#')
			if view[y][x] != scaffold {
				continue
			}

			// Check surrounding of a scaffold.
			if view[y-1][x] == scaffold && view[y+1][x] == scaffold && view[y][x-1] == scaffold && view[y][x+1] == scaffold {
				sum += x * y
			}
		}
	}
	fmt.Println("Sum:", sum)
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
	const ChannelSize = 16384

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
