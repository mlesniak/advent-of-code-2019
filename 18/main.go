package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"
)

func main() {
	view := load()

	sameArea := make(map[string]bool)

	paths := findPaths(view)
	fmt.Println("PATHS")
	for key, value := range paths {
		//fmt.Println(string(key), "=>")
		for _, value := range value {
			k := string(key) + string(value.key)
			sameArea[k] = true
			//fmt.Println("  -", value)
		}
	}

	// ------------------------------------------------------------------------------------------

	cs := []coordinate{}
	withInput(view, func(_x, _y, value int) {
		if value == '@' {
			cs = append(cs, coordinate{_x, _y})
		}
	})
	//fmt.Println(cs)

	// Für alle keys, die da sind, ich aber nicht erreichen konnte:
	// hinzufügen, erstmal kosten 0
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

			// Check if already inserted.
			if sameArea[string(key)+string(key2)] {
				continue
			}

			// Start code here.
			emptyFoundKeys := map[int]bool{}
			for _, sp := range cs {
				p := bfs(view, true, emptyFoundKeys, sp.x, sp.y, ik2)
				if p.length != -1 {
					//fmt.Println("Found path", p, "for", sp)
					delete(p.foundKeys, ik2)
					paths[ik] = append(paths[ik], p)
				}
			}

			//p1 := keys[ik]
			//p := bfs(view, true, emptyFoundKeys, p1.x, p1.y, ik2)
			//if p.length == -1 {
			//	start := keys[ik2]
			//	p2 := newPath(start)
			//	p2.key = ik2
			//	p2.foundKeys[ik] = true
			//	//p2.length = -1 // Mark as "special" path
			//	// Length is from starting point, i.e. only for beginning.
			//	paths[ik] = append(paths[ik], p2)
			//}
		}
	}

	// "Habe einen Pfad von jedem Key, der mich nichts kostet in alle anderen Bereiche, zu dem jeweiligen @"
	//for key, value := range paths {
	//	fmt.Println(string(key), "=>")
	//	for _, value := range value {
	//		fmt.Println("  -", value)
	//	}
	//}
	//if true {
	//	return
	//}
	// ------------------------------------------------------------------------------------------

	candidates := findInitialList(view)
	fmt.Println("\nINITIAL")
	fmt.Println(candidates)

	// Sorted list of already visited keys, format keyorder->minimal amount of steps
	cache := make(map[string]int)

	// d=100, f=102
	var minSolution *candidate
	i := 0
	for len(candidates) > 0 {
		//if i > 20 {
		//	return
		//}

		c := candidates[0]
		candidates = candidates[1:]

		// Add to history.
		sorted := string(c.key) + "." + SortString(c.path) + "." + string(c.path[len(c.path)-1])
		if cur, found := cache[sorted]; found {
			if cur > c.length {
				//fmt.Println("  Adding to cache:", sorted, "with length=", c.length)
				cache[sorted] = c.length
			}
		} else {
			//fmt.Println("  Adding to cache:", sorted, "with length=", c.length)
			cache[sorted] = c.length
		}

		i++
		if i%100000 == 0 {
			fmt.Print("\r", len(candidates), " ", candidates[0].length, " ", c.path, strings.Repeat(" ", 40))
			//wait()
		}

		if minSolution != nil && minSolution.length < c.length {
			continue
		}

		if len(c.foundKeys) == len(keys) {
			if minSolution == nil || minSolution.length > c.length {
				minSolution = &c
				fmt.Println("***", minSolution.length)
				fmt.Println(*minSolution)
				// TODO Prune candidate list?

				too := 0
				for _, value := range candidates {
					if value.length >= minSolution.length {
						too++
					}
				}
				fmt.Println(">:", too)
			}
			continue
		}

		fmt.Println("\nEXAM:", c)
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

			// Create whole new copy.
			nc := newCandidate
			nc.foundKeys = make(map[int]bool)
			for key, value := range c.foundKeys {
				if value {
					nc.foundKeys[key] = true
				}
			}
			for key, value := range newCandidate.foundKeys {
				if value {
					nc.foundKeys[key] = true
				}
			}
			nc.foundKeys[newCandidate.key] = true

			// Length computation.
			if sameArea[string(c.key)+string(nc.key)] {
				nc.length += c.length
			} else {
				// Walk path backwards, trying to find a value in the same area. If we have one,
				// use this distance.
				//fmt.Println("  CAND/external", nc)
				updLen := nc.length + c.length
			updateLoop:
				for i := len(c.path) - 1; i >= 0; i-- {
					if sameArea[string(nc.key)+string(c.path[i])] {
						//fmt.Println("  CAND.path", string(c.path[i]))
						// Search for path cost from lat element in the area to the new one.
						for _, v := range paths[int(c.path[i])] {
							if v.key == nc.key {
								updLen = v.length + c.length
								//fmt.Println("  updLen", updLen)
								break updateLoop
							}
						}
					}
				}
				nc.length = updLen
			}

			if minSolution != nil && minSolution.length < nc.length {
				continue
			}

			nc.path = c.path + string(nc.key)
			fmt.Println("  CAND", nc)

			// check cached value. If it is lower, ignore this candidate.
			//ncSorted := SortString(nc.path)
			ncSorted := string(nc.key) + "." + SortString(nc.path) + "." + string(nc.path[len(nc.path)-1])
			if limit, found := cache[ncSorted]; found {
				// Examine only if this is better.
				if nc.length < limit {
					//fmt.Println(" -- Examining, since better for", ncSorted)
					candidates = append([]candidate{nc}, candidates...)
				} else {
					//fmt.Println(" -- Better result for", ncSorted, "=", limit, "instead of", nc.length, ", ignoring")
				}
			} else {
				// Add if not cached
				candidates = append([]candidate{nc}, candidates...)
			}
		}
	}

	if minSolution != nil {
		fmt.Println(minSolution.length)
		fmt.Println(*minSolution)
	}
}

type cands []candidate

func (a cands) Len() int           { return len(a) }
func (a cands) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a cands) Less(i, j int) bool { return a[i].length < a[j].length }

func findInitialList(view [][]int) []candidate {

	cs := []coordinate{}
	withInput(view, func(_x, _y, value int) {
		if value == '@' {
			cs = append(cs, coordinate{_x, _y})
		}
	})

	res := []candidate{}
	for _, coord := range cs {
		candidates := findReachableKeys(view, nil, coord.x, coord.y)
		for idx := range candidates {
			candidates[idx].path = string(candidates[idx].key)
		}
		res = append(res, candidates...)
	}

	return res
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
				delete(p.foundKeys, ik2)
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
	fk := ""
	for key, _ := range p.foundKeys {
		fk += string(key)
	}

	dk := ""
	for key, _ := range p.doorsOnWay {
		dk += string(key)
	}
	//return fmt.Sprintf("%s/pos=<%v>/len=%d/foundKeys=%s.%d/doors=%s/path=%s", string(p.key), p.position, p.length, fk, len(p.foundKeys), dk, p.path)
	return fmt.Sprintf("%s/len=%d/foundKeys=%s.%d/doors=%s/path=%s", string(p.key), p.length, fk, len(p.foundKeys), dk, p.path)
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
