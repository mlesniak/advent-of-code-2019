package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type command struct {
	direction string
	steps     int
}

type wire []command
type plane [][]int

type coordinate struct {
	x int
	y int
}

func main() {
	wires := load()
	size := computeMaxSize(wires)

	// IDEA Instead of using a 2D array we could have also used a map with just the 'wired' coordinates to reduce
	// 		memory usage, since most cells are unused. Nevertheless, the current solution works for now.
	plane := allocatePlane(size)

	// map. intersection -> distance for each wire
	distances := make(map[coordinate]map[int]int)

	// Simulate wires.
	origin := initializeOrigin(plane)
	for id, wire := range wires {
		simulate(plane, origin, wire, id+1, distances)
	}
	// Restart first wire to find cuts with the following ones for correct distance computation.
	simulate(plane, origin, wires[0], 1, distances)

	//// Debugging ------------------------------
	//log.Printf("Resulting plane:\n")
	//for _, row := range plane {
	//	fmt.Printf("%v\n", row)
	//}
	//// Debugging ------------------------------

	for k, v := range distances {
		log.Printf("Cut %v\n", k)
		for id, steps := range v {
			log.Printf("  %d -> %d\n", id, steps)
		}
	}

	// Find all intersections, find minimum manhattan distance to origin.
	// Solution for first part.
	//cuts := findIntersections(plane)
	//findMinimumDistance(cuts, origin)

	// Iterate through all cuts and find lowest sum.
	minSum := math.MaxInt64
	for _, v := range distances {
		sum := 0
		for _, steps := range v {
			sum += steps
		}
		if minSum > sum {
			minSum = sum
		}
	}
	log.Printf("Minimal distance using steps %d\n", minSum)
}

func computeMaxSize(wires []wire) int {
	// Maximal distance in one direction.
	maxDistance := 0

	commands := []string{"R", "L", "D", "U"}
	for _, command := range commands {
		for _, wire := range wires {
			pathLength := 0
			// Single wire.
			for _, sc := range wire {
				if sc.direction == command {
					pathLength += sc.steps
				}
			}
			if maxDistance < pathLength {
				maxDistance = pathLength
			}
		}
	}

	// Origin is in the center of the rectangle.
	maxDistance *= 2 + 1

	log.Printf("Maximal distance %d\n", maxDistance)
	return maxDistance
}

func findMinimumDistance(cuts []coordinate, origin coordinate) {
	min := math.MaxInt64
	for _, cut := range cuts {
		mh := computeDistance(origin, cut)
		if mh < min {
			min = mh
		}
		log.Printf("%v -> %d\n", cut, mh)
	}
	log.Printf("Minimum distance %d\n", min)
}

func computeDistance(origin coordinate, cut coordinate) int {
	return int(math.Abs(float64(origin.x-cut.x)) + math.Abs(float64(origin.y-cut.y)))
}

func initializeOrigin(plane plane) coordinate {
	origin := coordinate{len(plane[0]) / 2, len(plane) / 2}
	plane[origin.y][origin.x] = 9
	log.Printf("Origin ox=%d, oy=%d\n", origin.x, origin.y)
	return origin
}

func findIntersections(plane plane) []coordinate {
	cs := make([]coordinate, 0)
	for row := range plane {
		for col := range plane {
			if plane[row][col] == 8 {
				cs = append(cs, coordinate{col, row})
			}
		}
	}

	log.Printf("Cuts at %v\n", cs)
	return cs
}

func simulate(plane plane, origin coordinate, wire wire, id int, cutSteps map[coordinate]map[int]int) {
	ox := origin.x
	oy := origin.y

	log.Printf("wire=%v\n", wire)
	steps := 0
	for _, command := range wire {
		switch command.direction {
		case "R":
			for i := 0; i < command.steps; i++ {
				ox++
				steps++
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this as a cut.
					plane[oy][ox] = 8
					addStepsToGlobalCuts(ox, oy, cutSteps, id, steps)
				} else {
					plane[oy][ox] = id
				}
			}
		case "U":
			for i := 0; i < command.steps; i++ {
				oy--
				steps++
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this as a cut.
					plane[oy][ox] = 8
					addStepsToGlobalCuts(ox, oy, cutSteps, id, steps)
				} else {
					plane[oy][ox] = id
				}
			}
		case "L":
			for i := 0; i < command.steps; i++ {
				ox--
				steps++
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this as a cut.
					plane[oy][ox] = 8
					addStepsToGlobalCuts(ox, oy, cutSteps, id, steps)
				} else {
					plane[oy][ox] = id
				}
			}
		case "D":
			for i := 0; i < command.steps; i++ {
				oy++
				steps++
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this as a cut.
					plane[oy][ox] = 8
					addStepsToGlobalCuts(ox, oy, cutSteps, id, steps)
				} else {
					plane[oy][ox] = id
				}
			}
		}
	}
}

func addStepsToGlobalCuts(ox int, oy int, cutSteps map[coordinate]map[int]int, id int, steps int) {
	fmt.Println("Found cut at ", ox, oy)

	// This is an intersection. Remember the steps for this wire, if not yet collected (checked above).
	coord := coordinate{ox, oy}
	cut := cutSteps[coord]
	if cut == nil {
		cut = make(map[int]int)
		cutSteps[coord] = cut
	}
	if _, ok := cut[id]; !ok {
		// Add steps since not been added before.
		cut[id] = steps
	}
}

func allocatePlane(size int) plane {
	// First idea: simulate on a virtual plane and do no fancy mathematics.
	plane := make([][]int, size)
	for y := range plane {
		plane[y] = make([]int, size)
	}
	log.Printf("Size=%dx%d\n", size, size)
	return plane
}

func load() []wire {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")

	// For testing, comment out this.
	//bytes = []byte("R5,U3\nU2,R6")

	var wires []wire
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		var wire wire
		commands := strings.Split(line, ",")
		for _, token := range commands {
			c := string(token[0])
			s, _ := strconv.Atoi(token[1:])
			wire = append(wire, command{c, s})
		}
		wires = append(wires, wire)
	}

	return wires
}
