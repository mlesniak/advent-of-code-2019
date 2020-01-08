package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type area [][]int

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

func main() {
	a := load()
	fmt.Println(a)
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
