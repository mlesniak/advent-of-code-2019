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

func (a levelArea) String() string {
	s := ""
	maxLevel := findMaximalLevel(a)

	for level := -maxLevel; level <= maxLevel; level++ {
		s += fmt.Sprintf("Depth:%d\n", level)
		s += a[level].String() + "\n\n"
	}

	return s
}

func (a levelArea) Next() levelArea {
	b := make(levelArea)

	// Find maximum level
	maxLevel := findMaximalLevel(a)

	// Compute new values for each of the existing levels.
	for level := -maxLevel; level <= maxLevel; level++ {
		// I don't like the call-API (parameters) of next, but this will suffice for now.
		b[level] = a[level].Next(level, a)
	}

	// Increase levels by -1 and +1 after computation is done.
	a[-maxLevel-1] = make(area)
	a[maxLevel+1] = make(area)
	b[-maxLevel-1] = a[-maxLevel-1].Next(-maxLevel-1, a)
	b[maxLevel+1] = a[maxLevel+1].Next(maxLevel+1, a)

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

func (a area) Next(level int, levelArea levelArea) area {
	b := make(area)

	// If we are near the subgrid, we have to consider levels below our own level.
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
			if row == 2 && col == 2 {
				s += "?"
				continue
			}

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

// Handle all special cases separately, do not try to be too clever.
func (a area) Neighbors(level int, levelArea levelArea, row int, col int) int {
	ns := 0

	// TODO Not considered: looking a level up (currently only below)

	// North
	nrow := row - 1
	ncol := col
	if nrow == 2 && ncol == 2 {
		// Look into level-1 (if it exists) and collect all values.
		down, levelExists := levelArea[level+1]
		if levelExists {
			for dcol := 0; dcol < 5; dcol++ {
				if down[point{4, dcol}] {
					ns++
				}
			}
		} else {
			// We have nothing for this level yet, hence everything is empty and we have no neighbours.
		}
	} else if a[point{nrow, ncol}] {
		ns++
	}

	// South
	nrow = row + 1
	ncol = col
	if nrow == 2 && ncol == 2 {
		// Look into level+1 (if it exists) and collect all values.
		down, levelExists := levelArea[level+1]
		if levelExists {
			for dcol := 0; dcol < 5; dcol++ {
				if down[point{0, dcol}] {
					ns++
				}
			}
		} else {
			// We have nothing for this level yet, hence everything is empty and we have no neighbours.
		}
	} else if a[point{nrow, ncol}] {
		ns++
	}

	// West
	nrow = row
	ncol = col - 1
	if nrow == 2 && ncol == 2 {
		// Look into level+1 (if it exists) and collect all values.
		down, levelExists := levelArea[level+1]
		if levelExists {
			for drow := 0; drow < 5; drow++ {
				if down[point{drow, 4}] {
					ns++
				}
			}
		} else {
			// We have nothing for this level yet, hence everything is empty and we have no neighbours.
		}
	} else if a[point{nrow, ncol}] {
		ns++
	}

	// East
	nrow = row
	ncol = col + 1
	if nrow == 2 && ncol == 2 {
		// Look into level+1 (if it exists) and collect all values.
		down, levelExists := levelArea[level+1]
		if levelExists {
			for drow := 0; drow < 5; drow++ {
				if down[point{drow, 0}] {
					ns++
				}
			}
		} else {
			// We have nothing for this level yet, hence everything is empty and we have no neighbours.
		}
	} else if a[point{nrow, ncol}] {
		ns++
	}

	return ns
}

func main() {
	a := load()

	la := make(levelArea)
	la[0] = a

	i := 0
	for {
		fmt.Println(la)
		la = la.Next()
		i++
		//if i > 1 {
		//	break
		//}
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
