package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

	// Implement simple BFS to find all possible key-to-key path lengths.
	sum := 0
	for key := 'a'; key <= 'z'; key++ {
		if _, found := keys[int(key)]; found == false {
			continue
		}
		for key2 := 'a'; key2 <= 'z'; key2++ {
			ik2 := int(key2)
			if _, found := keys[ik2]; found == false {
				continue
			}
			if key2 == key {
				continue
			}

			// Start code here.
			p1 := keys[int(key)]
			p := bfs(view, p1.x, p1.y, ik2)
			if p.length != -1 {
				fmt.Println(string(key), string(key2), p)
			}
		}
	}
	fmt.Println("Length", sum)
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

type path struct {
	position  coordinate
	length    int
	keysOnWay []int
}

func newPath(c coordinate) path {
	return path{
		position:  c,
		length:    0,
		keysOnWay: []int{},
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
		//if view[position.y][position.x] >= 'A' && view[position.y][position.x] <= 'Z' {
		//	continue
		//}

		// Check if found.
		if view[position.y][position.x] == key {
			return candidate
		}

		// Add found key on the way.
		foundKeys := []int{}
		foundKeys = append(foundKeys, candidate.keysOnWay...)
		if candidate.length > 0 && view[position.y][position.x] >= 'a' && view[position.y][position.x] <= 'z' {
			foundKeys = append(foundKeys, view[position.y][position.x])
		}

		// Generate new candidates.
		addCandidate(history, &candidates, position.x+1, position.y, candidate.length, foundKeys)
		addCandidate(history, &candidates, position.x-1, position.y, candidate.length, foundKeys)
		addCandidate(history, &candidates, position.x, position.y+1, candidate.length, foundKeys)
		addCandidate(history, &candidates, position.x, position.y-1, candidate.length, foundKeys)
	}

	return path{coordinate{0, 0}, -1, []int{}}
}

func addCandidate(history map[coordinate]bool, candidates *[]path, x int, y int, len int, foundKeys []int) {
	if history[coordinate{x, y}] == false {
		p := path{coordinate{x, y}, len + 1, foundKeys}
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
