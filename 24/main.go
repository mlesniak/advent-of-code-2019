package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type point struct {
	y int
	x int
}

// Non-empty fields are true.
type area map[point]bool

// Define the area for each level.
type levelArea map[int]area

func (a levelArea) Next() levelArea {
	b := make(levelArea)

	// Find maximum level
	maxLevel := findMaximalLevel(a)

	// Compute new values for each of the existing levels.
	for level := -maxLevel; level <= maxLevel; level++ {
		// I don't like the call-API (parameters) of next, but this will suffice for now.
		b[level] = a[level].Next(maxLevel, level, a)
	}

	// Increase levels by +1 and -1 after computation is done.
	return b
}

func findMaximalLevel(a levelArea) int {
	maxLevel := 0
	for level, _ := range a {
		if level > maxLevel {
			maxLevel = level
		}
	}
	return maxLevel
}

func (a area) Next(maxLevel, level int, levelArea levelArea) area {
	b := make(area)

	// If we are near the subgrid, we have to consider levels up and below of our level.

	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			// Only computation of neighbors has changed.
			ns := a.Neighbors(level, levelArea, row, col)

			// Rules are the same.
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

func (a area) Neighbors(level int, levelArea levelArea, row int, col int) int {
	ns := 0

	//if row == 2 && col == 2 {
	//todo()
	//}

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

	la := make(levelArea)
	la[0] = a

	for {
		fmt.Println(la[0])
		la = la.Next()
		wait()
	}

	//fmt.Println(a.Score())

	//history := make(map[int]bool)
	//for {
	//	a = a.Next()
	//	fmt.Println(a)
	//	score := a.Score()
	//	fmt.Println(score)
	//	if history[score] {
	//		fmt.Println(score)
	//		break
	//	}
	//	history[score] = true
	//}
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

func todo() {
	panic("TODO")
}