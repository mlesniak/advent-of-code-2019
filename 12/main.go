package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type planet struct {
	x, y, z    int
	vx, vy, vz int
}

func main() {
	planets := load()
	fmt.Println(planets)
}

func load() []planet {
	var planets []planet

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for idx, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimFunc(line, func(r rune) bool {
			if r == '>' || r == '<' {
				return true
			}
			return false
		})
		// Custom parser just for this format!
		parts := strings.Split(line, ",")
		fmt.Println(idx, parts)
		vec := make([]int, 3)
		for idx, coord := range parts {
			anum := strings.Split(coord, "=")[1]
			a, err := strconv.Atoi(anum)
			if err != nil {
				panic(err)
			}
			vec[idx] = a
		}
		planets = append(planets, planet{x: vec[0], y: vec[1], z: vec[2]})
	}

	return planets
}
