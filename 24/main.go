package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type area [][]int

func (a area) Score() int {
	sum := 0

	exp := 0
	for _, line := range a {
		for _, value := range line {
			if value == 1 {
				sum += int(math.Pow(float64(2), float64(exp)))
			}
			exp++
		}
	}

	return sum
}

func (a area) Next() area {
	// Init new area.
	b := make([][]int, len(a))
	for i, _ := range a[0] {
		b[i] = make([]int, len(a[0]))
	}

	for row, line := range a {
		for col, value := range line {
			ns := a.Neighbors(row, col)

			// Is there a bug?
			if value == 1 {
				if ns != 1 {
					b[row][col] = 0
				} else {
					b[row][col] = 1
				}
			}

			// Empty space.
			if value == 0 {
				if ns == 1 || ns == 2 {
					b[row][col] = 1
				} else {
					b[row][col] = 0
				}
			}
		}
	}

	return b
}

func (a area) String() string {
	s := ""

	for _, line := range a {
		for _, value := range line {
			switch value {
			case 0:
				s += "."
			case 1:
				s += "#"
			}
		}
		s += "\n"
	}

	return s
}

func (a area) Neighbors(row int, col int) int {
	ns := 0

	// Up
	if row > 0 && a[row-1][col] == 1 {
		ns++
	}
	// Down
	if row < len(a)-1 && a[row+1][col] == 1 {
		ns++
	}
	// Left
	if col > 0 && a[row][col-1] == 1 {
		ns++
	}
	// Right
	if col < len(a[row])-1 && a[row][col+1] == 1 {
		ns++
	}

	return ns
}

func main() {
	a := load()
	fmt.Println(a)
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

	a := make([][]int, len(lines))
	for row, line := range lines {
		a[row] = make([]int, len(line))
		for col, c := range line {
			switch c {
			case '.':
				a[row][col] = 0
			case '#':
				a[row][col] = 1
			}
		}
	}

	return a
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}
