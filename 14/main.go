package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type chemical struct {
	quantity int
	name     string
}

type equation struct {
	chemicals []chemical
	result    chemical
}

func (e chemical) String() string {
	return fmt.Sprintf("%d:%s", e.quantity, e.name)
}

func (e equation) String() string {
	return fmt.Sprintf("%v => %v", e.chemicals, e.result)
}

func main() {
	equations := load()
	showEquations(equations)

	storage := make(map[string]int)
	requirements := []chemical{equations["FUEL"].result}
	for len(requirements) > 0 {
		fmt.Println(strings.Repeat("-", 40))
		fmt.Println("Current requirements:", requirements)
		goal := requirements[0]
		fmt.Println("Current goal:", goal)
		requirements = requirements[1:]

		// Check if requirement depends only on an ore.
		dependents := equations[goal.name].chemicals
		if len(dependents) == 1 && dependents[0].name == "ORE" {
			fmt.Println("Found basic part", goal)
			storage[goal.name] += goal.quantity
			continue
		}

		// Use what we have in storage.
		if storage[goal.name] > 0 {
			capacity := storage[goal.name]

			if capacity >= goal.quantity {
				// Take everything out of storage.
				fmt.Println("For", goal.name, "taking EVERYTHING out of storage", capacity, "and continuing")
				storage[goal.name] = capacity - goal.quantity
				fmt.Println("No need to search for more.")
				continue
			}

			if capacity < goal.quantity {
				// Take something out of storage.
				fmt.Println("For", goal.name, "taking SOME out of storage", capacity, "and continuing")
				storage[goal.name] = 0
				goal.quantity = goal.quantity - capacity
			}
		}

		solution := findChemicals(equations, goal)
		fmt.Println("Solution for goal:", solution)
		if solution.result.quantity != goal.quantity {
			delta := solution.result.quantity - goal.quantity
			fmt.Println("Spare parts for", goal.name, ":", delta)
			storage[goal.name] = delta
		}
		requirements = append(requirements, solution.chemicals...)
	}

	fmt.Println("--- STORAGE ")
	ore := 0
	for key, value := range storage {
		if !isBasicPart(equations, key) {
			continue
		}
		o1 := findChemicals(equations, chemical{name: key, quantity: value})
		fmt.Println(key, value, o1, "=", o1.chemicals[0].quantity)
		ore += o1.chemicals[0].quantity
	}
	fmt.Println("Needed ore (13312):", ore)
}

func isBasicPart(equations map[string]equation, name string) bool {
	dependents := equations[name].chemicals
	return len(dependents) == 1 && dependents[0].name == "ORE"
}

func findChemicals(equations map[string]equation, goal chemical) equation {
	solution, found := equations[goal.name]
	if found == false {
		panic(fmt.Sprintf("No solution found: %v", goal))
	}

	sq := float64(solution.result.quantity)
	gq := float64(goal.quantity)
	if sq == gq {
		return solution
	}

	var adaptedEquation equation
	adaptedEquation.result = solution.result
	adaptedEquation.chemicals = make([]chemical, len(solution.chemicals))
	copy(adaptedEquation.chemicals, solution.chemicals)
	factor := int(math.Ceil(gq / sq))
	adaptedEquation.result.quantity = factor * solution.result.quantity
	for i, _ := range adaptedEquation.chemicals {
		adaptedEquation.chemicals[i].quantity *= factor
	}

	return adaptedEquation
}

func showEquations(equations map[string]equation) {
	log.Println("Parsed equations:")
	for k, e := range equations {
		log.Println(k, "=>", e.chemicals)
	}
}

func load() map[string]equation {
	equations := make(map[string]equation)

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		if len(strings.Trim(line, " ")) == 0 {
			continue
		}
		ps := strings.Split(line, "=>")
		result := parse(ps[1])
		var chemicals []chemical
		for _, c := range strings.Split(ps[0], ",") {
			chemicals = append(chemicals, parse(c))
		}
		equations[result.name] = equation{result: result, chemicals: chemicals}
	}

	return equations
}

func parse(component string) chemical {
	s := strings.Trim(component, " \t")
	ps := strings.Split(s, " ")
	q, err := strconv.Atoi(ps[0])
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
	return chemical{name: ps[1], quantity: q}
}
