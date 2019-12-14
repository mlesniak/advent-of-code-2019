package main

import (
	"fmt"
	"io/ioutil"
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
	//showEquations(equations)

	// Find all chemicals which need only ORE
	baseChemical := make(map[string]int)
	for _, e := range equations {
		if len(e.chemicals) == 1 && e.chemicals[0].name == "ORE" {
			baseChemical[e.result.name] = 0
		}
	}

	list := []chemical{equations["FUEL"].result}
	for len(list) > 0 {
		goal := list[0]
		list = list[1:]
		fmt.Println("\n---GOAL:", goal)

		if _, found := baseChemical[goal.name]; found {
			baseChemical[goal.name] += goal.quantity
			fmt.Println("BASE found", equations[goal.name].chemicals[0])
			continue
		}

		chemicals := findChemicals(equations, goal)
		list = append(list, chemicals.chemicals...)
	}

	ore := 0
	for key, needed := range baseChemical {
		e := equations[key]
		q := e.result.quantity
		factor := int(math.Ceil(float64(needed) / float64(q)))
		o := e.chemicals[0].quantity * factor
		ore += o
		fmt.Println(key, factor, o)
	}
	fmt.Println("ORE:", ore)
}

func findChemicals(equations map[string]equation, goal chemical) equation {
	solution, found := equations[goal.name]
	if found == false {
		panic(fmt.Sprintf("No solution found: %v", goal))
	}

	// TODO Multiply again.
	return solution
}

func showEquations(equations map[string]equation) {
	for k, e := range equations {
		fmt.Println(k, "=>", e.chemicals)
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
	s := strings.Trim(component, " ")
	ps := strings.Split(s, " ")
	q, err := strconv.Atoi(ps[0])
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
	return chemical{name: ps[1], quantity: q}
}

func oldVersion(baseChemical map[string]bool, equations map[string]equation) {
	//basics := []chemical{}
	//storage := make(map[string]int)
	//buildList := []chemical{chemical{1, "FUEL"}}
	//for len(buildList) > 0 {
	//	fmt.Println("\n- Step ------------------------")
	//	fmt.Println("* LIST", buildList)
	//
	//	goal := buildList[0]
	//	fmt.Println("Goal:", goal)
	//	buildList = buildList[1:]
	//
	//	_, ok := baseChemical[goal.name]
	//	if ok {
	//		// Base chemical.
	//		fmt.Println("BASE Chemical", goal)
	//		basics = append(basics, goal)
	//		continue
	//	}
	//
	//	//// If goal is ORE, we simply have it.
	//	//if goal.name == "ORE" {
	//	//	fmt.Println("ORE needed:", goal)
	//	//	ore += goal.quantity
	//	//	continue
	//	//}
	//
	//	// Check storage, if we have some chemicals left.
	//	//quantity, ok := storage[goal.name]
	//	//if ok {
	//	//	// Check if we have enough in storage. If yes, simply use it.
	//	//	if quantity >= goal.quantity {
	//	//		storage[goal.name] -= goal.quantity
	//	//		fmt.Println("Using storage", goal.quantity, "leaving", storage[goal.name])
	//	//		continue
	//	//	}
	//	//}
	//
	//	//solution := findChemicals(equations, goal)
	//	fmt.Println("Found", solution)
	//	// Update storage.
	//	leftOver := solution.result.quantity - goal.quantity
	//	if leftOver > 0 {
	//		storage[solution.result.name] += leftOver
	//		fmt.Println("Storing remaining", leftOver, "of", solution.result.name, "for now", storage[solution.result.name])
	//	}
	//
	//	// Update dependencies.
	//	buildList = append(buildList, solution.chemicals...)
	//}
	//fmt.Println("=====================================")
	//ores := make(map[string]int)
	//for _, value := range basics {
	//	ores[value.name] += value.quantity
	//}
	//fmt.Println(ores)
	//ore := 0
	////for key, value := range ores {
	////o1 := findChemicals(equations, chemical{name: key, quantity: value})
	////fmt.Println(key, value, o1)
	////ore += o1.chemicals[0].quantity
	////}
	//fmt.Println("*** ORE needed", ore)
}
