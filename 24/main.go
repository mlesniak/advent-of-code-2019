package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type point struct {
	y int
	x int
}

// Non-empty fields are true.
type area map[point]bool

func (a area) Score() int {
	sum := 0

	exp := 0
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			if a[point{row, col}] {
				sum += int(math.Pow(float64(2), float64(exp)))
			}
			exp++
		}
	}

	return sum
}

func (a area) Next() area {
	b := make(map[point]bool)

	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			ns := a.Neighbors(row, col)

			// Bug
			if a[point{row, col}] {
				if ns == 1 {
					b[point{row, col}] = true
				} else {
					b[point{row, col}] = false
				}
			}

			// Empty space.
			if !a[point{row, col}] {
				if ns == 1 || ns == 2 {
					b[point{row, col}] = true
				} else {
					b[point{row, col}] = false
				}
			}
		}
	}

	return b
}

func (a area) String() string {
	s := ""

	// Hard-coded size.
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			if a[point{row, col}] {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}

	return s
}

func (a area) Neighbors(row int, col int) int {
	ns := 0

	if a[point{row - 1, col}] {
		ns++
	}
	if a[point{row + 1, col}] {
		ns++
	}
	if a[point{row, col - 1}] {
		ns++
	}
	if a[point{row, col + 1}] {
		ns++
	}

	return ns
}

func main() {
	a := load()
	fmt.Println(a)
	//fmt.Println(a.Score())

	history := make(map[int]bool)
	for {
		a = a.Next()
		fmt.Println(a)
		score := a.Score()
		fmt.Println(score)
		if history[score] {
			fmt.Println(score)
			break
		}
		history[score] = true
	}
}

func load() area {
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	a := make(map[point]bool)
	for row, line := range lines {
		for col, c := range line {
			switch c {
			case '.':
				// Do nothing.
			case '#':
				a[point{row, col}] = true
			}
		}
	}

	return a
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}
