package main

import (
	"fmt"
	"io/ioutil"
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
	return fmt.Sprintf("%d %s", e.quantity, e.name)
}

func (e equation) String() string {
	return fmt.Sprintf("%v => %v", e.chemicals, e.result)
}

func main() {
	equations := load()
	//showEquations(equations)

	ore := 0
	storage := make(map[string]int)
	buildList := []chemical{chemical{1, "FUEL"}}
	for len(buildList) > 0 {
		fmt.Println("\n- Step ------------------------")
		fmt.Println("* LIST", buildList)

		goal := buildList[0]
		fmt.Println("Goal:", goal)
		buildList = buildList[1:]

		// If goal is ORE, we simply have it.
		if goal.name == "ORE" {
			fmt.Println("ORE needed:", goal)
			ore += goal.quantity
			continue
		}

		// Check storage, if we have some chemicals left.
		quantity, ok := storage[goal.name]
		if ok {
			// Check if we have enough in storage. If yes, simply use it.
			if quantity >= goal.quantity {
				storage[goal.name] -= goal.quantity
				fmt.Println("Using storage", goal.quantity, "leaving", storage[goal.name])
				continue
			}
		}

		solution := findChemicals(equations, goal)
		fmt.Println("Found", solution)
		// Update storage.
		leftOver := solution.result.quantity - goal.quantity
		if leftOver > 0 {
			storage[solution.result.name] += leftOver
			fmt.Println("Storing remaining", leftOver, "of", solution.result.name, "for now", storage[solution.result.name])
		}

		// Update dependencies.
		buildList = append(buildList, solution.chemicals...)
	}

	fmt.Println("\n\n*** ORE needed", ore)
}

// Brute force with added math.
func findChemicals(equations []equation, goal chemical) equation {
	var solution *equation
	for _, eq := range equations {
		if eq.result.name == goal.name {
			solution = &eq
			break
		}
	}
	if solution == nil {
		panic(fmt.Sprintf("No solution found: %v", goal))
	}

	cs := make([]chemical, len(solution.chemicals))
	copy(cs, solution.chemicals)
	// TODO Does not seem right...
	//for idx := range cs {
	//	cs[idx].quantity *= goal.quantity
	//}
	return *solution
}

func showEquations(equations []equation) {
	for _, e := range equations {
		fmt.Println(e)
	}
}

func load() []equation {
	var equations []equation

	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		ps := strings.Split(line, "=>")
		result := parse(ps[1])
		var chemicals []chemical
		for _, c := range strings.Split(ps[0], ",") {
			chemicals = append(chemicals, parse(c))
		}
		equations = append(equations, equation{result: result, chemicals: chemicals})
	}

	return equations
}

func parse(component string) chemical {
	s := strings.Trim(component, " ")
	ps := strings.Split(s, " ")
	q, err := strconv.Atoi(ps[0])
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
	return chemical{name: ps[1], quantity: q}
}
