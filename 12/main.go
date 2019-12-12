package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type planet struct {
	position vector
	velocity vector
}

func (p planet) String() string {
	return fmt.Sprintf("pos={%v} vel={%v}", p.position, p.velocity)
}

type vector struct {
	x, y, z int
}

func (v vector) String() string {
	return fmt.Sprintf("x=%d y=%d z=%d", v.x, v.y, v.z)
}

func (v vector) add(v2 vector) {
	v.x += v2.x
	v.y += v2.y
	v.z += v2.z
}

func main() {
	planets := load()
	fmt.Println(planets)

	const maxSteps = 10
	for step := 0; step <= maxSteps; step++ {
		// Status.
		fmt.Println("\nStep", step, strings.Repeat("-", 30))
		for _, p := range planets {
			fmt.Println(p)
		}

		// Update velocity for each planet
		velocities := computeVelocities(planets)

		// Apply velocity.
		for idx, v := range velocities {
			planets[idx].velocity = v
			planets[idx].position.add(planets[idx].velocity)
		}
	}
}

func computeVelocities(planets []planet) []vector {
	velocities := make([]vector, len(planets))
	return velocities
}

func load() []planet {
	var planets []planet

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
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
		vec := make([]int, 3)
		for idx, coord := range parts {
			anum := strings.Split(coord, "=")[1]
			a, err := strconv.Atoi(anum)
			if err != nil {
				panic(err)
			}
			vec[idx] = a
		}
		planets = append(planets, planet{position: vector{x: vec[0], y: vec[1], z: vec[2]}})
	}

	return planets
}
