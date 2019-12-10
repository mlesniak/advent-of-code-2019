package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type coordinate struct {
	x int
	y int
}

func main() {
	asteroids := load()
	fmt.Println(asteroids)

	// Determine maximum dimensions.
	maxX := 0
	maxY := 0
	for c := range asteroids {
		if c.x > maxX {
			maxX = c.x
		}
		if c.y > maxY {
			maxY = c.y
		}
	}

	// Compute hidden line of sight by computing x and y deltas.
	for asteroid := range asteroids {
		//if asteroid.x != 4 || asteroid.y != 4 {
		//	continue
		//}

		// Create a copy of the whole map to count visible asteroids.
		copy := make(map[coordinate]bool)
		for k, v := range asteroids {
			if k == asteroid {
				// Ignore self.
				continue
			}
			copy[k] = v
		}

		for candidate := range asteroids {
			if asteroid == candidate {
				continue
			}
			dx := candidate.x - asteroid.x
			dy := candidate.y - asteroid.y
			//fmt.Println(asteroid, candidate, dx, dy)

			// Remove all elements on the path from asteroid + delta
			px := asteroid.x
			py := asteroid.y
			for {
				//fmt.Println("  px=", px, "py=", py)
				// Out of bounds?
				if !(px >= 0 && px <= maxX && py >= 0 && py <= maxY) {
					break
				}
				// Delta-inducing asteroid?
				pc := coordinate{x: px, y: py}
				if pc != candidate {
					// Remove from map.
					delete(copy, pc)
				}

				// Next possible step.
				px += dx
				py += dy
			}

			// Count remaining asteroids in copy.
			// Remember largest map.
		}
		fmt.Println("***", asteroid, len(copy), " ||", copy)
	}
}

func load() map[coordinate]bool {
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	asteroids := make(map[coordinate]bool)
	for row, line := range lines {
		fmt.Println(row, line)
		for col, char := range line {
			if char == '#' {
				asteroids[coordinate{x: col, y: row}] = true
			}
		}
	}
	return asteroids
}
