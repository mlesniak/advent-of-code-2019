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
	data    [][]int
	portals map[point]point
	gates   map[string]point // start and end gate
}

func main() {
	data := load()

	//for key, value := range data.portals {
	//	fmt.Println(key, value)
	//}
	//if true {
	//	return
	//}

	start := data.gates["AA"]
	goal := data.gates["ZZ"]

	length := bfs(data, []path{{start, 0}}, goal)
	fmt.Println(length)
}

type path struct {
	position point
	length   int
}

func bfs(data maze, list []path, goal point) int {
	history := make(map[point]bool)
	for _, p := range list {
		history[p.position] = true
	}

	for len(list) > 0 {
		p := list[0]
		if p.position == goal {
			return p.length
		}

		list = list[1:]
		if p.length > 64 {
			continue
		}

		cs := getCandidates(data, p)
		for _, c := range cs {
			if history[c.position] == false {
				list = append(list, c)
				history[c.position] = true
			}
		}

		//wait()
	}

	// No path found.
	return -1
}

func getCandidates(data maze, p path) []path {
	// For a point, find all nearby possible ways.
	// TODO Add portals.
	result := []path{}
	view := data.data
	pos := p.position

	// Direct connection.
	if view[pos.y+1][pos.x] == '.' {
		np := path{point{pos.x, pos.y + 1}, p.length + 1}
		result = append(result, np)
	}
	if view[pos.y-1][pos.x] == '.' {
		np := path{point{pos.x, pos.y - 1}, p.length + 1}
		result = append(result, np)
	}
	if view[pos.y][pos.x+1] == '.' {
		np := path{point{pos.x + 1, pos.y}, p.length + 1}
		result = append(result, np)
	}
	if view[pos.y][pos.x-1] == '.' {
		np := path{point{pos.x - 1, pos.y}, p.length + 1}
		result = append(result, np)
	}

	// We could have portals on the points itself, i.e. not a step away, but this would be functionally incorrect
	// and we do not know what the second task is.
	pp := point{pos.x, pos.y + 1}
	if portal, found := data.portals[pp]; found {
		result = append(result, path{position: portal, length: p.length + 1})
	}
	pp = point{pos.x, pos.y - 1}
	if portal, found := data.portals[pp]; found {
		result = append(result, path{position: portal, length: p.length + 1})
	}
	pp = point{pos.x + 1, pos.y}
	if portal, found := data.portals[pp]; found {
		result = append(result, path{position: portal, length: p.length + 1})
	}
	pp = point{pos.x - 1, pos.y}
	if portal, found := data.portals[pp]; found {
		result = append(result, path{position: portal, length: p.length + 1})
	}

	return result
}

func load() maze {
	data := make([][]int, 0)

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	// Determine maximum column length.
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	for y, row := range lines {
		colMem := make([]int, maxLen)
		data = append(data, colMem)
		for x, col := range row {
			data[y][x] = int(col)
		}
	}

	directionGates := computeDirectionGates(data)
	portals := computePortalExits(directionGates)
	gates := computeStartEndPoints(directionGates)

	return maze{data: data, portals: portals, gates: gates}
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
