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

func (p point) String() string {
	return fmt.Sprintf("[%d,%d]", p.y, p.x)
}

type portal struct {
	point       // Goal
	inside bool // Is the source of the portal on the inside?
}

type maze struct {
	data    [][]int
	portals map[point]portal
	gates   map[string]point // start and end gate
}

func main() {
	fmt.Println("Start")
	data := load()
	//for key, value := range debugPortals {
	//    fmt.Println(key, value)
	//}
	//for key, value := range data.portals {
	//    fmt.Println(key, value, value.inside)
	//}
	//wait()

	start := data.gates["AA"]
	goal := data.gates["ZZ"]
	fmt.Println("AA", start)
	fmt.Println("ZZ", goal)

	length := bfs(data, []path{{start, 0, 0}}, goal)
	fmt.Println(length)
}

type path struct {
	position point
	length   int
	level    int
}

func (p path) String() string {
	return fmt.Sprintf("%v/length=%d/level=%d", p.position, p.length, p.level)
}

func bfs(data maze, list []path, goal point) int {
	type historyPoint struct {
		point
		level int
	}
	history := make(map[historyPoint]bool)
	for _, p := range list {
		history[historyPoint{point: p.position, level: 0}] = true
	}

	for len(list) > 0 {
		p := list[0]
		fmt.Println("\nExploring", p)
		if name, found := debugPortals[p.position]; found {
			fmt.Println("* before portal", name)
		}

		if p.position == goal && p.level == 0 {
			return p.length
		}

		list = list[1:]
		if p.length > 512 {
			continue
		}

		cs := getCandidates(data, p)
		for _, c := range cs {
			if history[historyPoint{c.position, c.level}] == false {
				fmt.Println("+ candidate", c)
				list = append(list, c)
				history[historyPoint{c.position, c.level}] = true
			}
		}

		//wait()
	}

	// No path found.
	return -1
}

func getCandidates(data maze, p path) []path {
	// For a point, find all nearby possible ways.
	result := []path{}
	view := data.data
	pos := p.position

	// Direct connection.
	if view[pos.y+1][pos.x] == '.' {
		np := path{point{pos.x, pos.y + 1}, p.length + 1, p.level}
		result = append(result, np)
	}
	if view[pos.y-1][pos.x] == '.' {
		np := path{point{pos.x, pos.y - 1}, p.length + 1, p.level}
		result = append(result, np)
	}
	if view[pos.y][pos.x+1] == '.' {
		np := path{point{pos.x + 1, pos.y}, p.length + 1, p.level}
		result = append(result, np)
	}
	if view[pos.y][pos.x-1] == '.' {
		np := path{point{pos.x - 1, pos.y}, p.length + 1, p.level}
		result = append(result, np)
	}

	// We could have portals on the points itself, i.e. not a step away, but this would be functionally incorrect
	// and we do not know what the second task is.
	pp := point{pos.x, pos.y + 1}
	if portal, found := data.portals[pp]; found && (p.level >= 1 || portal.inside) {
		if pt, f := debugPortals[portal.point]; f {
			fmt.Println("\t\t\t\t\t\tUsing portal", pt)
		}
		delta := -1
		if portal.inside {
			delta = +1
		}
		result = append(result, path{position: portal.point, length: p.length + 1, level: p.level + delta})
	}
	pp = point{pos.x, pos.y - 1}
	if portal, found := data.portals[pp]; found && (p.level >= 1 || portal.inside) {
		if pt, f := debugPortals[portal.point]; f {
			fmt.Println("\t\t\t\t\t\tUsing portal", pt)
		}
		delta := -1
		if portal.inside {
			delta = +1
		}
		result = append(result, path{position: portal.point, length: p.length + 1, level: p.level + delta})
	}
	pp = point{pos.x + 1, pos.y}
	if portal, found := data.portals[pp]; found && (p.level >= 1 || portal.inside) {
		if pt, f := debugPortals[portal.point]; f {
			fmt.Println("\t\t\t\t\t\tUsing portal", pt)
		}
		delta := -1
		if portal.inside {
			delta = +1
		}
		result = append(result, path{position: portal.point, length: p.length + 1, level: p.level + delta})
	}
	pp = point{pos.x - 1, pos.y}
	if portal, found := data.portals[pp]; found && (p.level >= 1 || portal.inside) {
		if pt, f := debugPortals[portal.point]; f {
			fmt.Println("\t\t\t\t\t\tUsing portal", pt)
		}
		delta := -1
		if portal.inside {
			delta = +1
		}
		result = append(result, path{position: portal.point, length: p.length + 1, level: p.level + delta})
	}

	return result
}

func newPortal(maze [][]int, x int, y int) portal {
	// Check if we are on the outside
	inside := true
	if x <= 1 || x-2 >= len(maze[y]) || y <= 1 || y-2 >= len(maze) {
		inside = false
	}
	return portal{point{x, y}, inside}
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
	portals := computePortalExits(data, directionGates)
	gates := computeStartEndPoints(data, directionGates)

	return maze{data: data, portals: portals, gates: gates}
}

func computeStartEndPoints(maze [][]int, directionGates map[string][]pointDir) map[string]point {
	gates := make(map[string]point)
	for key, value := range directionGates {
		if len(value) > 1 {
			key = key + ".1"
		}

		p1 := value[0]
		gates[key] = adapt(maze, p1.point, p1.point, p1.orientation).point
	}
	return gates
}

var debugPortals = make(map[point]string)

func computePortalExits(maze [][]int, directionGates map[string][]pointDir) map[point]portal {
	portals := make(map[point]portal)
	for key, value := range directionGates {
		// Ignore start and end gate.
		if len(value) < 2 {
			continue
		}

		p1 := value[0]
		p2 := value[1]
		portals[p1.point] = adapt(maze, p1.point, p2.point, p2.orientation)
		portals[p2.point] = adapt(maze, p2.point, p1.point, p1.orientation)

		s := "/OUT"
		if portals[p1.point].inside {
			s = "/IN"
		}
		debugPortals[portals[p2.point].point] = key + s
		s = "/OUT"
		if portals[p2.point].inside {
			s = "/IN"
		}
		debugPortals[portals[p1.point].point] = key + s
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

func adapt(maze [][]int, src point, p point, orientation int) portal {
	inside := true
	if src.x <= 1 || src.x >= len(maze[src.y])-3 || src.y <= 1 || src.y >= len(maze)-2 {
		inside = false
	}

	switch orientation {
	case 1:
		return portal{point{p.x, p.y - 1}, inside}
	case 2:
		return portal{point{p.x, p.y + 1}, inside}
	case 3:
		return portal{point{p.x - 1, p.y}, inside}
	case 4:
		return portal{point{p.x + 1, p.y}, inside}
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
