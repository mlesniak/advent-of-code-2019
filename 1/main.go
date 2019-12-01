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
		sum += compute(mass)
	}
	fmt.Println(sum)

	//fmt.Println(compute(12))
	//fmt.Println(compute(14))
	//fmt.Println(compute(1969))
	//fmt.Println(compute(100756))
}

func compute(n int) int {
	return n/3 - 2
}
