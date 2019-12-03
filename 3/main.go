package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type command struct {
	direction string
	steps     int
}

type wire []command
type plane [][]int

func main() {
	// TODO We could add up all single R, L, U, and Ds to find possible maximal dimensions?
	const size = 10
	wires := load()
	plane := allocatePlane(size)

	// Simulate wires.
	for _, wire := range wires {
		simulate(plane, wire)
	}

	// Find intersections.
}

func simulate(plane plane, wire wire) {
	ox := len(plane[0]) / 2
	oy := len(plane) / 2
	log.Printf("Origin ox=%d, oy=%d\n", ox, oy)
	plane[oy][ox] = 9

	log.Printf("wire=%v\n", wire)
	// "Paint" plane with 1, if a cable is at position,
	// and with a 2 if it intersects an existing cable.
	for _, command := range wire {
		switch command.direction {
		case "R":
			for i := 0; i < command.steps; i++ {
				ox++
				plane[oy][ox] = plane[oy][ox] + 1
			}
		case "U":
			for i := 0; i < command.steps; i++ {
				oy--
				plane[oy][ox] = plane[oy][ox] + 1
			}
		case "L":
			for i := 0; i < command.steps; i++ {
				ox--
				plane[oy][ox] = plane[oy][ox] + 1
			}
		case "D":
			for i := 0; i < command.steps; i++ {
				oy++
				plane[oy][ox] = plane[oy][ox] + 1
			}
		}
	}

	log.Printf("Resulting wire:\n")
	for _, row := range plane {
		fmt.Printf("%v\n", row)
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
	//bytes, _ := ioutil.ReadFile("input.txt")
	// For testing, see https://adventofcode.com/2019/day/2.
	bytes := []byte("R8,U5,L5,D3")

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
