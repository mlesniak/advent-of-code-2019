package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	sum := 0
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		mass, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// https://adventofcode.com/2019/day/1
		// Specifically, to find the fuel required for a module, take its mass,
		// divide by three, round down, and subtract 2.
		fuel := compute(mass)
		sum += fuel

		// https://adventofcode.com/2019/day/1#part2
		sum += fuelCompute(fuel)
	}
	fmt.Println(sum)

	//fmt.Println(compute(12))
	//fmt.Println(compute(14))
	//fmt.Println(compute(1969))
	//fmt.Println(compute(100756))

	//fmt.Println(fuelCompute(654))
}

func fuelCompute(mass int) int {
	additionalFuel := compute(mass)
	if additionalFuel <= 0 {
		return 0
	}
	return additionalFuel + fuelCompute(additionalFuel)
}

func compute(n int) int {
	return n/3 - 2
}
