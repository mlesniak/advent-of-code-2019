package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type point struct {
	x int
	y int
}

type maze struct {
	data   [][]int
	portal map[point]point
	gates  map[string][]point
}

func main() {
	load()

}

func load() maze {
	data := make([][]int, 0)

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for y, row := range lines {
		colMem := make([]int, len(row))
		data = append(data, colMem)
		for x, col := range row {
			data[y][x] = int(col)
		}
	}

	type pointDir struct {
		point
		//    1
		//   3 4
		//    2
		orientation int
	}

	// Find gates.
	gates := make(map[string][]pointDir)
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			if data[y][x] >= 'A' && data[y][x] <= 'Z' {
				// Look around if a corresponding dot can be found.
				var id string
				var p pointDir
				if y > 0 && data[y-1][x] == '.' {
					id = SortString(string(data[y+1][x]) + string(data[y][x]))
					p = pointDir{point{x, y}, 1}
				}
				if y < len(data)-1 && data[y+1][x] == '.' {
					id = SortString(string(data[y-1][x]) + string(data[y][x]))
					p = pointDir{point{x, y}, 2}
				}
				if x < len(data[0])-1 && data[y][x+1] == '.' {
					id = SortString(string(data[y][x-1]) + string(data[y][x]))
					p = pointDir{point{x, y}, 4}
				}
				if x > 0 && data[y][x-1] == '.' {
					id = SortString(string(data[y][x+1]) + string(data[y][x]))
					p = pointDir{point{x, y}, 3}
				}
				if id != "" {
					gates[id] = append(gates[id], p)
				}
			}
		}
	}

	for key, value := range gates {
		// Ignore start and end gate.
		if len(value) > 1 {

		}
		fmt.Println(key, value)
	}
	portals := make(map[point]point)
	return maze{data: data, portal: portals}
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func SortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}
