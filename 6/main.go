package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type node struct {
	name   string
	parent *node
}

func main() {
	defs := load()

	refs := make(map[string]*node)
	for _, def := range defs {
		src := getNodeOrCreate(refs, def.source)
		dest := getNodeOrCreate(refs, def.destination)
		dest.parent = src
	}

	//steps := 0
	//for _, def := range refs {
	//	steps += countLinkage(def)
	//}
	//fmt.Println(steps)

	pathYou := getPath(refs, refs["YOU"])
	pathSan := getPath(refs, refs["SAN"])
	youSteps := 0
loop:
	for _, v := range pathYou {
		youSteps++
		sanSteps := 0
		for _, w := range pathSan {
			sanSteps++
			if v == w {
				// Common ancestor found.
				youSteps += sanSteps
				break loop
			}
		}
	}
	// Ignore steps to the parent of YOU and SAN.
	youSteps -= 2
	fmt.Println(youSteps)
}

func getPath(nodes map[string]*node, node *node) []string {
	if node.name == "COM" {
		return []string{""}
	}

	path := getPath(nodes, node.parent)
	return append([]string{node.parent.name}, path...)
}

func countLinkage(def *node) int {
	if def.parent == nil {
		return 0
	}
	if def.parent.name == "COM" {
		return 1
	}

	return 1 + countLinkage(def.parent)
}

func getNodeOrCreate(refs map[string]*node, name string) *node {
	val, ok := refs[name]
	if !ok {
		refs[name] = &node{
			name:   name,
			parent: nil,
		}
		val = refs[name]
	}

	return val
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
