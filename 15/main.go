package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MemorySize = 1000000

// Directions.
const (
	north = 1
	south = 2
	west  = 3
	east  = 4
)

// Replies.
const (
	wall    = 0
	ok      = 1
	success = 2
	drone   = 3
)

func fromDirection(direction int) string {
	switch direction {
	case 1:
		return "north"
	case 2:
		return "south"
	case 3:
		return "west"
	case 4:
		return "east"
	}

	panic(direction)
}

func fromReply(reply int) string {
	switch reply {
	case 0:
		return "wall"
	case 1:
		return "ok"
	case 2:
		return "found"
	}

	panic(reply)
}

type step struct {
	direction int
	tried     map[int]bool
}

func (s step) String() string {
	keys := make([]string, 0)
	for key, _ := range s.tried {
		keys = append(keys, fromDirection(key))
	}

	return fmt.Sprintf("%s%v", fromDirection(s.direction), keys)
}

type path []step

func debug(a ...interface{}) {
	fmt.Println(a...)
}

func main() {
	memory := load()
	in := newChannel()
	out := newChannel()

	height := 45
	width := 45

	canvas := make([][]int, height)
	for row := range canvas {
		canvas[row] = make([]int, width)
		for i := 0; i < width; i++ {
			canvas[row][i] = -1
		}
	}
	y := len(canvas) / 2
	x := len(canvas[0]) / 2

	// A path is an array of steps.
	startStep := step{north, map[int]bool{north: true}}
	path := []step{startStep}
	go func() {
		for {
			canvas[y][x] = drone
			canvas[len(canvas)/2][len(canvas)/2] = 10

			//fmt.Print("?")
			//bufio.NewReader(os.Stdin).ReadLine()
			//debug("\nPath", path)
			fmt.Println(len(path))
			//if len(path) > 20 {
			//	panic("")
			//}

			// Walk a step into the given direction.
			direction := path[len(path)-1].direction
			debug("Walking", fromDirection(direction))
			in <- direction

			reply := <-out
			debug("Received reply", fromReply(reply))
			switch reply {
			case wall:
				switch direction {
				case north:
					canvas[y-1][x] = wall
				case south:
					canvas[y+1][x] = wall
				case east:
					canvas[y][x+1] = wall
				case west:
					canvas[y][x-1] = wall
				}
				path = backtrack(path)
				debug("Backtracked...")
			case ok:
				canvas[y][x] = ok
				switch direction {
				case north:
					y--
				case south:
					y++
				case east:
					x++
				case west:
					x--
				}
				canvas[y][x] = drone
				// Add next path step in same direction and ignore way back.
				s := path[len(path)-1]
				newStep := step{s.direction, map[int]bool{reverse(s.direction): true}}
				path = append(path, newStep)
			case success:
				canvas[y][x] = ok
				switch direction {
				case north:
					canvas[y-1][x] = ok
					y--
				case south:
					canvas[y+1][x] = ok
					y++
				case east:
					canvas[y][x+1] = ok
					x++
				case west:
					canvas[y][x-1] = ok
					x--
				}
				canvas[y][x] = drone
				debug("***", x, y, len(path))
				panic("found")
			}

			paintCanvas(canvas)
		}
	}()
	compute("15", memory, in, out)

}

func reverse(direction int) int {
	switch direction {
	case north:
		return south
	case south:
		return north
	case east:
		return west
	case west:
		return east
	}

	panic(direction)
}

func backtrack(path path) path {
	for {
		if len(path) == 0 {
			panic("No way found.")
		}

		// If possible, choose another option for last element in the path.
		last := len(path) - 1
		step := &path[last]
		if len(step.tried) == 4 {
			//	// Remove last element and restart with previous one.
			debug("Backtracking...")
			path = path[:len(path)-1]
		} else {
			// Choose another direction which was not chosen before.
			for i := 1; i <= 4; i++ {
				_, found := step.tried[i]
				if !found {
					step.tried[i] = true
					step.direction = i
					return path
				}
			}
		}
	}
}

func paintCanvas(canvas [][]int) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("\033[2J")

	for row := range canvas {
		for col := range canvas[row] {
			var c string
			switch canvas[row][col] {
			case -1:
				c = "."
			case ok:
				c = " "
			case wall:
				c = "#"
			case drone:
				c = "D"
			case 10:
				c = "X"
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
	time.Sleep(time.Millisecond * 10)
}

func newChannel() chan int {
	channelSize := 16384
	return make(chan int, channelSize)
}

func compute(name string, memory []int, in chan int, out chan int) {
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

func load() []int {
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
	return memory
}
