package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	view := load()

	//fmt.Println("PATHS")
	paths := findPaths(view)
	//for key, value := range paths {
	//	fmt.Println(string(key), "=>")
	//	for _, value := range value {
	//		fmt.Println("  -", value)
	//	}
	//}

	keys := findKeys(view)

	candidates := findInitialList(view)
	//fmt.Println("\nINITIAL")
	//fmt.Println(candidates)

	var minSolution *candidate

	i := 0
	for len(candidates) > 0 {
		c := candidates[0]
		candidates = candidates[1:]

		i++
		if i%1000 == 0 {
			log.Println(len(candidates))
			log.Println(c.path)
			//wait()
		}

		if minSolution != nil && minSolution.length < c.length {
			// Ignore longer paths
			continue
		}

		if len(c.foundKeys) == len(keys) {
			if minSolution == nil || minSolution.length > c.length {
				minSolution = &c
				fmt.Println("***", minSolution.length)
				fmt.Println(*minSolution)
			}
			//
			//fmt.Println("Solution", c)
			//fmt.Println("***", c.length)
			continue
		}

		//fmt.Println("\nEXAM:", c)
		//cs := findReachableKeys(view, c.foundKeys, c.position.x, c.position.y)
		cs := paths[c.key]

	nextCandidate:
		for _, newCandidate := range cs {
			// If already in keys, ignore.
			if c.foundKeys[newCandidate.key] {
				continue
			}

			// Check if for all doors, a key exists.
			for doorKey, _ := range newCandidate.doorsOnWay {
				if !c.foundKeys[doorKey] {
					continue nextCandidate
				}
			}

			// Update length
			newCandidate.length += c.length
			newCandidate.path = c.path + string(newCandidate.key)
			//fmt.Println("  CAND", newCandidate)
			if minSolution != nil && minSolution.length < c.length {
				// Ignore longer paths
				continue
			}
			candidates = append([]candidate{newCandidate}, candidates...)
			//candidates = append(candidates, newCandidate)
		}
	}

	if minSolution != nil {
		fmt.Println(minSolution.length)
		fmt.Println(*minSolution)
	}
}

func findInitialList(view [][]int) []candidate {
	var x, y int
	withInput(view, func(_x, _y, value int) {
		if value == '@' {
			x = _x
			y = _y
		}
	})
	candidates := findReachableKeys(view, nil, x, y)
	for idx := range candidates {
		candidates[idx].path = string(candidates[idx].key)
	}
	return candidates
}

func findPaths(view [][]int) map[int][]candidate {
	paths := map[int][]candidate{}
	keys := findKeys(view)
	for key := 'a'; key <= 'z'; key++ {
		ik := int(key)
		if _, found := keys[ik]; found == false {
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
			p1 := keys[ik]
			emptyFoundKeys := map[int]bool{}
			p := bfs(view, true, emptyFoundKeys, p1.x, p1.y, ik2)
			if p.length != -1 {
				paths[ik] = append(paths[ik], p)
			}
		}
	}
	return paths
}

func findReachableKeys(view [][]int, foundKeys map[int]bool, x int, y int) []candidate {
	candidates := []candidate{}
	for key := 'a'; key <= 'z'; key++ {
		// If already found, ignore.
		if foundKeys[int(key)] {
			continue
		}

		ik := int(key)
		p := bfs(view, false, foundKeys, x, y, ik)
		if p.length > 0 {
			candidates = append(candidates, p)
		}
	}
	return candidates
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

type mbool map[int]bool

func (p mbool) String() string {
	s := ""
	for key, value := range map[int]bool(p) {
		if value {
			s = s + string(key) + ","
		}
	}
	s = s[:len(s)-1]
	return s
}

func (p candidate) String() string {
	return fmt.Sprintf("%s/pos=<%v>/len=%d/foundKeys=%v/doors=%v/path=%s", string(p.key), p.position, p.length, p.foundKeys, p.doorsOnWay, p.path)
}

type candidate struct {
	position   coordinate
	path       string
	length     int
	key        int
	foundKeys  map[int]bool
	doorsOnWay map[int]bool
}

func newPath(c coordinate) candidate {
	return candidate{
		position:   c,
		length:     0,
		foundKeys:  map[int]bool{},
		doorsOnWay: map[int]bool{},
	}
}

func bfs(view [][]int, ignoreDoors bool, _foundKeys map[int]bool, x int, y int, key int) candidate {
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
		// TODO Comment this when finished...
		if view[position.y][position.x] >= 'A' && view[position.y][position.x] <= 'Z' {
			lower := int(unicode.ToLower(rune(view[position.y][position.x])))
			if !ignoreDoors && !candidate.foundKeys[lower] {
				continue
			}
			candidate.doorsOnWay[lower] = true
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

		dw := map[int]bool{}
		for k, v := range c.doorsOnWay {
			dw[k] = v
		}

		p := candidate{
			position:   coordinate{x, y},
			length:     c.length + 1,
			key:        key,
			foundKeys:  fk,
			doorsOnWay: dw,
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
