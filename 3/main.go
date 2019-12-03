package main

import (
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
	// TODO We could add up all single R, L, U, and Ds to find possible maximal dimensions?
	wires := load()
	size := computeMaxSize(wires)
	plane := allocatePlane(size)

	// Simulate wires.
	origin := initializeOrigin(plane)
	for id, wire := range wires {
		simulate(plane, origin, wire, id+1)
	}

	// Find all intersections, find minimum manhattan distance to origin.
	cuts := findIntersections(plane)
	findMinimumDistance(cuts, origin)
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
			if plane[row][col] == -1 {
				cs = append(cs, coordinate{col, row})
			}
		}
	}

	log.Printf("Cuts at %v\n", cs)
	return cs
}

func simulate(plane plane, origin coordinate, wire wire, id int) {
	ox := origin.x
	oy := origin.y

	log.Printf("wire=%v\n", wire)
	// "Paint" plane with 1, if a cable is at position,
	// and with a 2 if it intersects an existing cable.
	for _, command := range wire {
		switch command.direction {
		case "R":
			for i := 0; i < command.steps; i++ {
				ox++
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this.
					plane[oy][ox] = -1
				} else {
					plane[oy][ox] = id
				}
			}
		case "U":
			for i := 0; i < command.steps; i++ {
				oy--
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this.
					plane[oy][ox] = -1
				} else {
					plane[oy][ox] = id
				}
			}
		case "L":
			for i := 0; i < command.steps; i++ {
				ox--
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this.
					plane[oy][ox] = -1
				} else {
					plane[oy][ox] = id
				}
			}
		case "D":
			for i := 0; i < command.steps; i++ {
				oy++
				if plane[oy][ox] != 0 && plane[oy][ox] != id {
					// There is already a wire. Mark this.
					plane[oy][ox] = -1
				} else {
					plane[oy][ox] = id
				}
			}
		}
	}

	//log.Printf("Resulting wire:\n")
	//for _, row := range plane {
	//	fmt.Printf("%v\n", row)
	//}
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
