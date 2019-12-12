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

func (v *vector) add(v2 vector) {
	v.x += v2.x
	v.y += v2.y
	v.z += v2.z
}

func main() {
	planets := load()
	fmt.Println(planets)

	var history [][]planet

	// We have a lot of time...
	const maxSteps = 2772 + 1
	//const maxSteps = 100

loop:
	for step := 0; step <= maxSteps; step++ {
		// Remember this state.
		c := make([]planet, len(planets))
		copy(c, planets)
		history = append(history, c)

		// Status.
		showStatus(step, planets)

		if step == maxSteps {
			// Do not compute velocity for last step.
			continue
		}

		// Update velocity for each planet
		velocities := computeVelocities(planets)

		// Apply velocity.
		for idx, v := range velocities {
			planets[idx].velocity = v
			planets[idx].position.add(planets[idx].velocity)
		}

		// Did we already saw this values?
		for _, previous := range history {
			if compare(planets, previous) {
				showStatus(step+1, planets)
				fmt.Println("\n\n*** Repeated state after steps:", step+1)
				break loop
			}
		}
	}
}

func showStatus(step int, planets []planet) {
	fmt.Println("\nStep", step, strings.Repeat("-", 30))
	for _, p := range planets {
		fmt.Println(p)
	}
}

func computeEnergy() {
	// Compute energy.
	//energy := 0.0
	//for idx := range planets {
	//	p := planets[idx]
	//	pose := math.Abs(float64(p.position.x)) + math.Abs(float64(p.position.y)) + math.Abs(float64(p.position.z))
	//	vele := math.Abs(float64(p.velocity.x)) + math.Abs(float64(p.velocity.y)) + math.Abs(float64(p.velocity.z))
	//	total := pose * vele
	//	//fmt.Println(idx, pose, vele, total)
	//	energy += total
	//}
	//fmt.Println("\nENERGY", energy)
}

// There must be a better way?
func compare(current []planet, previous []planet) bool {
	for idx := range current {
		if current[idx] != previous[idx] {
			return false
		}
	}
	return true
}

func computeVelocities(planets []planet) []vector {
	velocities := make([]vector, len(planets))
	for idx := range velocities {
		velocities[idx] = planets[idx].velocity
	}

	// Woohooo, O(n^2) ...
	for i := 0; i < len(planets); i++ {
		for j := 0; j < len(planets); j++ {
			if i == j {
				continue
			}
			p1 := planets[i]
			p2 := planets[j]

			velocities[i].x += computeGravity(p1.position.x, p2.position.x)
			velocities[i].y += computeGravity(p1.position.y, p2.position.y)
			velocities[i].z += computeGravity(p1.position.z, p2.position.z)
		}
	}

	return velocities
}

func computeGravity(x1 int, x2 int) int {
	switch {
	case x1 == x2:
		return 0
	case x1 > x2:
		return -1
	case x1 < x2:
		return +1
	}

	panic("Not reachable")
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
