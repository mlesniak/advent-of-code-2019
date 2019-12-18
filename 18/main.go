package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func main() {
	view := load()

	// Find starting point.
	var x, y int
	withInput(view, func(_x, _y, value int) {
		if value == '@' {
			x = _x
			y = _y
		}
	})

	keys := findKeys(view)
	doors := findDoors(view)

	// Find initial list of reachable nodes for search.
	candidates := findReachableKeys(view, x, y)
	fmt.Println("Initial candidates", candidates)

	// Idea:DFS sorted by current path length?
	i := 0
	for len(candidates) > 0 {
		i++
		if i > 10 {
			break
		}
		// TODO Each candidate needs its own view!
		fmt.Println("\n\nCandidates", candidates)
		candidate := candidates[0]
		candidates = candidates[1:]
		fmt.Println("Examining", string(candidate.key), candidate)

		// Copy view since we will open doors and remove keys.
		cp := make([][]int, len(view))
		for row := range view {
			cp[row] = make([]int, len(view[0]))
			for col := range view[row] {
				cp[row][col] = view[row][col]
			}
		}

		// Remove key and open door for this key.
		keyCoord := keys[candidate.key]
		cp[keyCoord.y][keyCoord.x] = '.'
		doorCoord := doors[int(unicode.ToUpper(rune(candidate.key)))]
		cp[doorCoord.y][doorCoord.x] = '.'

		// Find now reachable keys and add them to the list.
		cs := findReachableKeys(cp, keyCoord.x, keyCoord.y)
		for _, c := range cs {
			c.view = cp
		}
		candidates = append(candidates, cs...)

		// If none can be found, ... TODO
	}
}

func findReachableKeys(view [][]int, x int, y int) []path {
	candidates := []path{}
	for key := 'a'; key <= 'c'; key++ {
		ik := int(key)
		p := bfs(view, x, y, ik)
		if p.length > 0 {
			candidates = append(candidates, p)
		}
	}
	return candidates
}

func generateMatrix() {
	//// Implement simple BFS to find all possible key-to-key path lengths.
	//sum := 0
	//for key := 'a'; key <= 'z'; key++ {
	//	if _, found := keys[int(key)]; found == false {
	//		continue
	//	}
	//	for key2 := 'a'; key2 <= 'z'; key2++ {
	//		ik2 := int(key2)
	//		if _, found := keys[ik2]; found == false {
	//			continue
	//		}
	//		if key2 == key {
	//			continue
	//		}
	//
	//		// Start code here.
	//		p1 := keys[int(key)]
	//		p := bfs(view, p1.x, p1.y, ik2)
	//		if p.length != -1 {
	//			fmt.Println(string(key), string(key2), p)
	//		}
	//	}
	//}
	//fmt.Println("Length", sum)
}

func findDoors(view [][]int) map[int]coordinate {
	doors := make(map[int]coordinate)
	for c := 'A'; c <= 'Z'; c++ {
		withInput(view, func(x, y, value int) {
			ic := int(c)
			if value == ic {
				doors[ic] = coordinate{x, y}
			}
		})
	}
	return doors
}

func findKeys(view [][]int) map[int]coordinate {
	keys := make(map[int]coordinate)
	for c := 'a'; c <= 'z'; c++ {
		withInput(view, func(x, y, value int) {
			ic := int(c)
			if value == ic {
				keys[ic] = coordinate{x, y}
			}
		})
	}
	return keys
}

func simpleSearchFromAtoZ() {
	//fmt.Println("Searching", string(key), ", starting at", x, y)
	//result := bfs(view, x, y, int(key))
	//fmt.Println("->", result)
	//
	//// Unlock corresponding door.
	//doorPosition := doors[int(unicode.ToUpper(rune(key)))]
	//view[doorPosition.y][doorPosition.x] = '.'
	//
	//// Start search for next key.
	//x = result.position.x
	//y = result.position.y
	//sum += result.length
}

func (p path) String() string {
	return fmt.Sprintf("%s/pos=<%v>/len=%d/visited=%v", string(p.key), p.position, p.length, p.visited)
}

type path struct {
	position  coordinate
	length    int
	key       int
	view      [][]int
	foundKeys []int
	visited   map[int]bool
}

func newPath(c coordinate) path {
	return path{
		position:  c,
		length:    0,
		view:      [][]int{},
		foundKeys: []int{},
		visited:   map[int]bool{},
	}
}

func bfs(view [][]int, x int, y int, key int) path {
	start := coordinate{x, y}
	candidates := []path{newPath(start)}
	history := make(map[coordinate]bool)

	for len(candidates) > 0 {
		//wait()
		//fmt.Println(candidates)
		candidate := candidates[0]
		position := candidate.position
		//fmt.Println("- Looking at", candidate)
		candidates = candidates[1:]

		// Ignore already visited.
		_, visited := history[position]
		if visited {
			continue
		}
		history[position] = true

		// Ignore out-of-fields.
		if position.x < 0 || position.x > len(view[0]) || position.y < 0 || position.y > len(view) {
			continue
		}
		// Ignore walls.
		if view[position.y][position.x] == '#' {
			continue
		}
		// Ignore doors.
		if view[position.y][position.x] >= 'A' && view[position.y][position.x] <= 'Z' {
			continue
		}

		// Check if found.
		if view[position.y][position.x] == key {
			return candidate
		}

		// Generate new candidates.
		addCandidate(history, &candidates, key, position.x+1, position.y, candidate)
		addCandidate(history, &candidates, key, position.x-1, position.y, candidate)
		addCandidate(history, &candidates, key, position.x, position.y+1, candidate)
		addCandidate(history, &candidates, key, position.x, position.y-1, candidate)
	}

	p := newPath(coordinate{0, 0})
	p.length = -1
	return p
}

func addCandidate(history map[coordinate]bool, candidates *[]path, key int, x int, y int, candidate path) {
	if history[coordinate{x, y}] == false {
		p := path{
			position:  coordinate{x, y},
			length:    candidate.length + 1,
			key:       key,
			foundKeys: candidate.foundKeys,
			view:      candidate.view,
		}
		*candidates = append(*candidates, p)
	}
}

type coordinate struct {
	x, y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("%d/%d", c.x, c.y)
}

type direction int

const (
	up = iota
	down
	left
	right
)

const minDirection = up
const maxDirection = right

func (d direction) String() string {
	switch d {
	case up:
		return "U"
	case down:
		return "D"
	case left:
		return "L"
	case right:
		return "R"
	}

	panic("Unknown direction")
}

func oppositeDirection(d direction) direction {
	switch d {
	case up:
		return down
	case down:
		return up
	case left:
		return right
	case right:
		return left
	}

	panic("Unknown direction")
}

func fromDelta(dx int, dy int) direction {
	switch {
	case dx == 0 && dy == 1:
		return down
	case dx == 0 && dy == -1:
		return up
	case dx == 1 && dy == 0:
		return right
	case dx == -1 && dy == 0:
		return left
	}

	panic("Unknown dx/dy")
}

func withInput(input [][]int, f func(x, y, value int)) {
	for row := range input {
		for col := range input[row] {
			f(col, row, input[row][col])
		}
	}
}

func load() [][]int {
	input := make([][]int, 0)

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for y, row := range lines {
		colMem := make([]int, len(row))
		input = append(input, colMem)
		for x, col := range row {
			input[y][x] = int(col)
		}
	}

	return input
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}
