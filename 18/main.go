package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
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
	//doors := findDoors(view)

	// Find initial list of reachable nodes for search.
	candidates := findReachableKeys(view, nil, x, y)
	fmt.Println("Initial candidates", candidates)

	minSolution := math.MaxInt64
	for len(candidates) > 0 {
		//fmt.Print("\r", len(candidates))
		//fmt.Println("\n------------------------------------\nCandidates", candidates)
		candidate := candidates[0]
		candidates = candidates[1:]
		fmt.Println("Examining", string(candidate.key), candidate)

		// Check if this is a solution.
		if len(keys) == len(candidate.foundKeys) {
			if minSolution > candidate.length {
				minSolution = candidate.length
			}
			////fmt.Println("*** Solution with length=", candidate.length)
			//solutions = append(solutions, candidate.length)
			continue
		}

		// Find now reachable keys and add them to the list.
		keyPosition := keys[candidate.key]
		cs := findReachableKeys(view, candidate.foundKeys, keyPosition.x, keyPosition.y)
		for idx, _ := range cs {
			fmt.Println("-> Candidate:", cs[idx])
			cs[idx].length = cs[idx].length + candidate.length
		}
		// DFS instead of BFS.
		candidates = append(cs, candidates...)
	}

	fmt.Println()
	fmt.Println(minSolution)
}

func findReachableKeys(view [][]int, foundKeys map[int]bool, x int, y int) []candidate {
	candidates := []candidate{}
	for key := 'a'; key <= 'z'; key++ {
		// If already found, ignore.
		if foundKeys[int(key)] {
			continue
		}

		ik := int(key)
		p := bfs(view, foundKeys, x, y, ik)
		if p.length > 0 {
			candidates = append(candidates, p)
		}
	}
	return candidates
}

func generateMatrix() {
	//// Implement simple BFS to find all possible key-to-key candidate lengths.
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

func (p candidate) String() string {
	return fmt.Sprintf("%s/pos=<%v>/len=%d/foundKeys=%v", string(p.key), p.position, p.length, p.foundKeys)
}

type candidate struct {
	position  coordinate
	length    int
	key       int
	foundKeys map[int]bool
}

func newPath(c coordinate) candidate {
	return candidate{
		position:  c,
		length:    0,
		foundKeys: map[int]bool{},
	}
}

func bfs(view [][]int, _foundKeys map[int]bool, x int, y int, key int) candidate {
	start := coordinate{x, y}
	path := newPath(start)
	for k, v := range _foundKeys {
		path.foundKeys[k] = v
	}
	candidates := []candidate{path}
	history := make(map[coordinate]bool)

	for len(candidates) > 0 {
		candidate := candidates[0]
		position := candidate.position
		candidates = candidates[1:]

		// Ignore already foundKeys.
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
		// Ignore doors for found keys, stop otherwise.
		if view[position.y][position.x] >= 'A' && view[position.y][position.x] <= 'Z' {
			lower := int(unicode.ToLower(rune(view[position.y][position.x])))
			if !candidate.foundKeys[lower] {
				continue
			}
		}
		// Add keys on the way to the list of found keys.
		if view[position.y][position.x] >= 'a' && view[position.y][position.x] <= 'z' {
			lower := int(unicode.ToLower(rune(view[position.y][position.x])))
			candidate.foundKeys[lower] = true
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

func addCandidate(history map[coordinate]bool, candidates *[]candidate, key int, x int, y int, c candidate) {
	if history[coordinate{x, y}] == false {
		// Copy foundKeys
		fk := map[int]bool{}
		for k, v := range c.foundKeys {
			fk[k] = v
		}

		p := candidate{
			position:  coordinate{x, y},
			length:    c.length + 1,
			key:       key,
			foundKeys: fk,
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
