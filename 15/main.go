package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func debug(a ...interface{}) {
	fmt.Println(a...)
}

func main() {
	memory, in, out := load()

	size := 50
	ship := make([][]int, size)
	for row := range ship {
		ship[row] = make([]int, size)
		for key, _ := range ship[row] {
			ship[row][key] = -1
		}
	}
	x := len(ship[0]) / 2
	y := len(ship) / 2
	ship[y][x] = 3
	//fmt.Println("begin", x, y)

	maxLen := 650
	if len(os.Args) > 1 {
		m, err := strconv.Atoi(os.Args[1])
		maxLen = m
		if err != nil {
			panic(err)
		}
	}
	go func() {
		backtrack(ship, x, y, maxLen, in, out, 0, nil)
		drawShip(ship)
		computeOxygen(ship)
		os.Exit(1) // HACK
	}()
	compute(memory, in, out)
}

func computeOxygen(ship [][]int) {
	minutes := 0

	for {
		todo := findNumberOfPlacesToFill(ship)
		if todo == 0 {
			break
		}
		minutes++
		drawShip(ship)
		fmt.Println(todo)
		//wait()
		// Status
		fmt.Println("Places to fill", todo)
		filled := findNumberOfPlacesFilled(ship)
		fmt.Println("Places filled", filled)

		// Copy data to prevent racing.
		cp := make([][]int, len(ship))
		for i := range ship {
			cp[i] = make([]int, len(ship[i]))
			copy(cp[i], ship[i])
		}

		// For all oxygen places
		withShip(ship, func(x, y, val int) {
			if val == 2 {
				// Oxygen found.
				// Look at surrounding places and check if a non-filled space exists.
				if ship[y+1][x] == 1 {
					cp[y+1][x] = 2
				}
				if ship[y-1][x] == 1 {
					cp[y-1][x] = 2
				}
				if ship[y][x+1] == 1 {
					cp[y][x+1] = 2
				}
				if ship[y][x-1] == 1 {
					cp[y][x-1] = 2
				}
			}
		})
		ship = cp
	}

	fmt.Println("Filling took minutes=", minutes)
}

func findNumberOfPlacesFilled(ship [][]int) int {
	filled := 0
	withShip(ship, func(x, y, val int) {
		if val == 2 {
			filled++
		}
	})
	return filled
}

func findNumberOfPlacesToFill(ship [][]int) int {
	todo := 0
	withShip(ship, func(x, y, val int) {
		if val == 1 {
			todo++
		}
	})
	return todo
}

func withShip(ship [][]int, f func(x, y, value int)) {
	for row := range ship {
		for col := range ship[row] {
			f(col, row, ship[row][col])
		}
	}
}

func backtrack(ship [][]int, x int, y int, maxLen int, in chan int, out chan int, length int, path []int) {
	if maxLen == length {
		//fmt.Println("<-")
		return
	}

	//wait()
	//drawShip(ship)
	//fmt.Println(path)

	//             1
	//           3 3 3 3 3 3
	//			   2
	//

	// Directions:
	//       1
	//     3   4
	//       2
	directions := []int{1, 2, 3, 4}
	for _, direction := range directions {
		// Do not choose the direct reversal since we would be staying at the previous step.
		//fmt.Println("@path=", path)
		if len(path) > 0 && opposite(path[len(path)-1]) == direction {
			//fmt.Println("not choosing reversal")
			continue
		}

		in <- direction
		//fmt.Println("?", direction)
		reply := <-out
		//fmt.Println(">", reply)
		switch reply {
		case 0: // Wall
			switch direction {
			case 1:
				ship[y-1][x] = 0
			case 2:
				ship[y+1][x] = 0
			case 3:
				ship[y][x-1] = 0
			case 4:
				ship[y][x+1] = 0
			}
			continue
		case 1: // OK
			ship[y][x] = 1
			switch direction {
			case 1:
				y--
			case 2:
				y++
			case 3:
				x--
			case 4:
				x++
			}
			ship[y][x] = 3

			backtrack(ship, x, y, maxLen, in, out, length+1, append(path, direction))
			ship[y][x] = 1
			opp := opposite(direction)
			switch opp {
			case 1:
				y--
			case 2:
				y++
			case 3:
				x--
			case 4:
				x++
			}
			in <- opp
			r := <-out
			if r != 1 {
				panic("Should not happen!")
			}
		case 2: // Energy source
			switch direction {
			case 1:
				y--
			case 2:
				y++
			case 3:
				x--
			case 4:
				x++
			}
			ship[y][x] = 2

			backtrack(ship, x, y, maxLen, in, out, length+1, append(path, direction))
			opp := opposite(direction)
			switch opp {
			case 1:
				y--
			case 2:
				y++
			case 3:
				x--
			case 4:
				x++
			}
			in <- opp
			r := <-out
			if r != 1 {
				panic("Should not happen!")
			}
		}
	}
}

func opposite(dir int) int {
	if dir == 1 {
		return 2
	}
	if dir == 2 {
		return 1
	}
	if dir == 3 {
		return 4
	}
	if dir == 4 {
		return 3
	}

	panic("Unsupported argument:" + string(dir))
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()

	//time.Sleep(time.Millisecond * 25)
}

func compute(memory []int, in chan int, out chan int) {
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
			return
		default:
			panic("Unknown opcode:" + strconv.Itoa(opcode))
		}
	}
}

const MemorySize = 1000000
const ChannelSize = 16384

func load() ([]int, chan int, chan int) {
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
	return memory, in, out
}

func drawShip(ship [][]int) {
	x := len(ship[0]) / 2
	y := len(ship) / 2

	//os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
	for row := range ship {
		for col := range ship[row] {
			var c string
			switch ship[row][col] {
			case -1:
				c = "_"
			case 0:
				c = "#"
			case 1:
				c = "."
			case 2:
				c = "O️"
			case 3:
				c = "S"
			}
			if col == x && row == y {
				c = "X"
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}
