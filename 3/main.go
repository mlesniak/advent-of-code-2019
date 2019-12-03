package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type command struct {
	direction string
	steps     int
}

type wire []command

func main() {
	wires := load()
	fmt.Println(wires)

}

func load() []wire {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	// For testing, see https://adventofcode.com/2019/day/2.
	//bytes := []byte("1,9,10,3,2,3,11,0,99,30,40,50")

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
