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

//func (a area) Next() area {
//	// Init new area.
//	b := make([][]int, len(a))
//	for i, _ := range a[0] {
//		b[i] = make([]int, len(a[0]))
//	}
//
//	for row, line := range a {
//		for col, value := range line {
//			ns := a.Neighbors(row, col)
//
//			// Is there a bug?
//			if value == 1 {
//				if ns != 1 {
//					b[row][col] = 0
//				} else {
//					b[row][col] = 1
//				}
//			}
//
//			// Empty space.
//			if value == 0 {
//				if ns == 1 || ns == 2 {
//					b[row][col] = 1
//				} else {
//					b[row][col] = 0
//				}
//			}
//		}
//	}
//
//	return b
//}

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

//func (a area) Neighbors(row int, col int) int {
//	ns := 0
//
//	// Up
//	if row > 0 && a[row-1][col] == 1 {
//		ns++
//	}
//	// Down
//	if row < len(a)-1 && a[row+1][col] == 1 {
//		ns++
//	}
//	// Left
//	if col > 0 && a[row][col-1] == 1 {
//		ns++
//	}
//	// Right
//	if col < len(a[row])-1 && a[row][col+1] == 1 {
//		ns++
//	}
//
//	return ns
//}

func main() {
	a := load()
	fmt.Println(a)
	fmt.Println(a.Score())

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
