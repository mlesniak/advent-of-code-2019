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
	gates  map[string]point // start and end gate
}

func main() {
	data := load()

	start := data.gates["AA"]
	goal := data.gates["ZZ"]

	length := bfs([]point{start}, goal)
	fmt.Println(length)
}

func bfs(list []point, goal point) int {
	for len(list) > 0 {
		p := list[0]
		list = list[1:]

		fmt.Println(p)
	}

	// No path found.
	return -1
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

	directionGates := computeDirectionGates(data)
	portals := computePortalExits(directionGates)
	gates := computeStartEndPoints(directionGates)

	return maze{data: data, portal: portals, gates: gates}
}

func computeStartEndPoints(directionGates map[string][]pointDir) map[string]point {
	gates := make(map[string]point)
	for key, value := range directionGates {
		if len(value) > 1 {
			continue
		}

		p1 := value[0]
		gates[key] = adapt(p1.point, p1.orientation)
	}
	return gates
}

func computePortalExits(directionGates map[string][]pointDir) map[point]point {
	portals := make(map[point]point)
	for _, value := range directionGates {
		// Ignore start and end gate.
		if len(value) < 2 {
			continue
		}

		p1 := value[0]
		p2 := value[1]
		portals[p1.point] = adapt(p2.point, p2.orientation)
		portals[p2.point] = adapt(p1.point, p1.orientation)
	}
	return portals
}

type pointDir struct {
	point
	//    1
	//   3 4
	//    2
	orientation int
}

func computeDirectionGates(data [][]int) map[string][]pointDir {
	directionGates := make(map[string][]pointDir)
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
					directionGates[id] = append(directionGates[id], p)
				}
			}
		}
	}
	return directionGates
}

func adapt(p point, orientation int) point {
	switch orientation {
	case 1:
		return point{p.x, p.y - 1}
	case 2:
		return point{p.x, p.y + 1}
	case 3:
		return point{p.x - 1, p.y}
	case 4:
		return point{p.x + 1, p.y}
	}

	panic(fmt.Sprintf("Unknown orientation: %d", orientation))
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
