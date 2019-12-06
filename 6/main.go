package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type node struct {
	name     string
	children []*node
}

func main() {
	defs := load()

	refs := make(map[string]node)
	for _, def := range defs {
		// add map to single nodes
		_, ok := refs[def.source]
		if !ok {
			refs[def.source] = node{
				name:     def.source,
				children: []*node{},
			}
		}
	}

	fmt.Println(refs)
}

type def struct {
	source      string
	destination string
}

func load() []def {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	// For testing, see https://adventofcode.com/2019/day/2.
	//bytes := []byte("1,9,10,3,2,3,11,0,99,30,40,50")
	lines := strings.Split(string(bytes), "\n")

	var defs []def
	for _, line := range lines {
		parts := strings.Split(line, ")")
		defs = append(defs, def{parts[0], parts[1]})
	}
	return defs
}
