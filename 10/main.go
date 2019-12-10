package main

import (
	"fmt"
	"io/ioutil"
	"math"
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
		//block := coordinate{4,2}
		//if asteroid != block {
		//	continue
		//}
		//fmt.Println("*** Debugging block", block)

		// Create a copy of the whole map to count visible asteroids.
		copy := make(map[coordinate]bool)
		for k, v := range asteroids {
			if k == asteroid {
				// Ignore self.
				continue
			}
			copy[k] = v
		}

		// Remember set of slopes.
		slopes := make(map[float64]bool)

		for candidate := range asteroids {
			if asteroid == candidate {
				continue
			}
			slope := computeSlope(asteroid, candidate)
			//dx := candidate.x - asteroid.x
			//dy := candidate.y - asteroid.y
			//slope := float32(dx) / float32(dy)

			// If the slope exists already, we have a blocking asteroid.
			// NOTE: currently we do not check that we get the nearest one,
			// which is not necessary for this problem.
			_, slopeExisting := slopes[slope]
			if slopeExisting {
				//fmt.Println("Ignoring", candidate, "with slope", slope)
				delete(copy, candidate)
				continue
			}

			// Otherwise, remember slope.
			slopes[slope] = true
			//fmt.Println(asteroid, candidate, dx, dy, "slope:", slope)

			// Remove all elements on the path from asteroid + delta
			//px := asteroid.x
			//py := asteroid.y
			//for {
			//	fmt.Println("  px=", px, "py=", py)
			//	// Out of bounds?
			//	if !(px >= 0 && px <= maxX && py >= 0 && py <= maxY) {
			//		break
			//	}
			//	// Delta-inducing asteroid?
			//	pc := coordinate{x: px, y: py}
			//	if pc != candidate {
			//		// Remove from map.
			//		delete(copy, pc)
			//	}
			//
			//	// Next possible step.
			//	px += dx
			//	py += dy
			//}
		}

		fmt.Println("***", asteroid, len(copy), " ||", copy)
	}
}

func computeSlope(origin coordinate, point coordinate) float64 {
	dx := float64(point.x - origin.x)
	dy := float64(point.y - origin.y)

	if dx < 0 && dy < 0 {
		return math.Atan(dy/dx) - math.Pi
	}
	if dx == 0 && dy < 0 {
		return -math.Pi / 2
	}
	if dx > 0 {
		return math.Atan(dy / dx)
	}
	if dx == 0 && dy > 0 {
		return math.Pi / 2
	}
	if dx < 0 && dy >= 0 {
		return math.Atan(dy/dx) + math.Pi
	}

	panic("Case not found")
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
