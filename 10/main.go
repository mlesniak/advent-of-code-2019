//
// NOTE ---------------------------------------------------------------------------------------------
// This code is ugly and ... works but it not optimized or refactored in any way. Let's hope that we
// do not have another puzzle of this kind. In this case I'd probably rewrite the whole thing.
//
// ¯\_(ツ)_/¯
//
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

type polar struct {
	distance float64
	rho      float64
}

func (p polar) String() string {
	return fmt.Sprintf("dist=%g, degree=%g", p.distance, p.rho)
}

func main() {
	asteroids := load()
	//part1(asteroids)

	// Compute polar coordinates for each other asteroid.
	polars := make(map[coordinate]polar)
	//station := coordinate{8, 3}
	station := coordinate{31, 20}
	fmt.Println("Station at", station)
	for asteroid := range asteroids {
		if station == asteroid {
			continue
		}
		p := computePolar(station, asteroid)
		polars[asteroid] = p
		fmt.Println(asteroid, p)
	}

	// Until all asteroids have been removed.
	lastRho := 0 - math.SmallestNonzeroFloat64
	i := 1
	for len(polars) > 0 {
		//fmt.Println("------------------------------------------------------------------------", i)
		// Find smallest rho larger than the last one. Special case at the beginning: rho = 0
		minRho := math.MaxFloat64
		for _, p := range polars {
			if p.rho > lastRho && p.rho < minRho {
				minRho = p.rho
			}
		}
		lastRho = minRho
		//fmt.Println("lastRho/minRho", lastRho)

		// Find smallest distance (and thus: asteroid) for the given rho value.
		minDist := math.MaxFloat64
		// Check why it does not work with *coordinate?
		target := coordinate{-1, -1}
		for c, p := range polars {
			if p.rho != minRho {
				continue
			}
			if p.distance < minDist {
				minDist = p.distance
				target = c
			}
		}
		nope := coordinate{-1, -1}
		if target == nope {
			lastRho = 0 - math.SmallestNonzeroFloat64
			continue
			//panic("No target found, resetting")
		}

		fmt.Println(i, "destroyed asteroid is at", target)
		delete(polars, target)
		i++

		//if i >= 25 {
		//	break
		//	fmt.Println("Breakpoint")
		//}
	}
}

func part1(asteroids map[coordinate]bool) {
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
	max := -1
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
		}

		if len(copy) > max {
			fmt.Println(asteroid, len(copy))
			max = len(copy)
		}
		//fmt.Println("***", asteroid, len(copy), " ||", copy)
	}
	fmt.Println(max)
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

// Use degree instead of polar coordinates for easier understanding.
func computePolar(origin coordinate, point coordinate) polar {
	rho := computeSlope(origin, point) * 180 / math.Pi

	dx := float64(point.x - origin.x)
	dy := float64(point.y - origin.y)
	dist := math.Sqrt(dx*dx + dy*dy)
	f := rho + 90.0
	if f < 0 {
		f = f + 360
	}
	return polar{dist, f}
}

func load() map[coordinate]bool {
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	asteroids := make(map[coordinate]bool)
	for row, line := range lines {
		//fmt.Println(row, line)
		for col, char := range line {
			if char == '#' {
				asteroids[coordinate{x: col, y: row}] = true
			}
		}
	}
	return asteroids
}
