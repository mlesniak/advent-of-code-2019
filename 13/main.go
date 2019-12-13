package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MemorySize = 1000000

func main() {
	// Idea: Stay under the ball by following its left and right movement.

	score := -1
	memory, in, out, canvas := initializeGame(&score)
	//handleManualInput(in)
	in <- 0
	compute("memory", memory, in, out)
	paintCanvas(canvas, score)
	println(score)
}

func nextInput(input []int) bool {
	idx := 0
	for {
		if idx == len(input) {
			return false
		}
		input[idx]++
		if input[idx] == 2 {
			input[idx] = -1
			idx++
		} else {
			return true
		}
	}
}

func initializeGame(score *int) ([]int, chan int, chan int, [][]int) {
	paddlePosition := 20
	scoreVisited := false

	// Create canvas.
	canvas := make([][]int, 24)
	for row := range canvas {
		canvas[row] = make([]int, 44)
	}
	memory := load()
	in := newChannel()
	out := newChannel()
	var prevT int
	go func() {
		for {
			// Do not paint anything if we simply simulate the game.
			if scoreVisited && prevT != 0 {
				//paintCanvas(canvas, *score)
			}

			x := <-out
			y := <-out
			t := <-out
			prevT = t

			// Track ball movement.
			// Ball position changed.
			var nextInput int
			if t == 4 && scoreVisited {
				currentBallX := x
				if currentBallX > paddlePosition {
					nextInput = 1
				}
				if currentBallX < paddlePosition {
					nextInput = -1
				}
				if currentBallX == paddlePosition {
					nextInput = 0
				}
				fmt.Println("-CurrentBall", x, "paddlePosition", paddlePosition, "-> input", nextInput)
				//fmt.Println("-Sending input", nextInput)
				in <- nextInput
			}
			// Update paddle position.
			if t == 3 {
				paddlePosition = x
			}

			// Score handling.
			if x == -1 {
				*score = t
				scoreVisited = true
			} else {
				canvas[y][x] = t
			}
		}
	}()
	// Allow free games.
	memory[0] = 2
	return memory, in, out, canvas
}

func countBlocks() {
	//// Count blocks.
	//blocks := 0
	//forCanvas(canvas, func(x, y, val int) {
	//	if val == 2 {
	//		blocks++
	//	}
	//})
	//println(blocks)
}

func handleManualInput(in chan int) chan int {
	go func() {
		for {
			r := bufio.NewReader(os.Stdin)
			line, _, err := r.ReadLine()
			if err != nil {
				panic(err)
			}
			command, _ := strconv.Atoi(string(line))
			in <- command
			//time.Sleep(time.Millisecond)
		}
	}()
	return in
}

func forCanvas(canvas [][]int, f func(int, int, int)) {
	for row := range canvas {
		for col := range canvas[row] {
			f(col, row, canvas[row][col])
		}
	}
}

func paintCanvas(canvas [][]int, score int) {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println(strings.Repeat("-", 80))
	for row := range canvas {
		for col := range canvas[row] {
			var c string
			switch canvas[row][col] {
			case 0:
				c = " "
			case 1:
				c = "#"
			case 2:
				c = "*"
			case 3:
				c = "-"
			case 4:
				c = "O️"
			}
			fmt.Print(c)
		}
		fmt.Println()
	}

	fmt.Println("SCORE", score)
	time.Sleep(time.Millisecond * 24)
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
