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
				// TODO Update something?
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
		o1 := findChemicals(equations, chemical{name: key, quantity: value})
		fmt.Println(key, value, o1)
		ore += o1.chemicals[0].quantity
	}
	fmt.Println("Needed ore (13312):", ore)
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

	//fmt.Println("*** solution", adaptedEquation, ", goal", goal)
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

// TODO Wenn es keine Basics sind, wird nicht reused.

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

func version2(equations map[string]equation) {
	//// Find all chemicals which need only ORE
	//baseChemical := make(map[string]int)
	//for _, e := range equations {
	//	if len(e.chemicals) == 1 && e.chemicals[0].name == "ORE" {
	//		fmt.Println("BASE element", e.result.name)
	//		baseChemical[e.result.name] = 0
	//	}
	//}
	//// Store unused quantities.
	//storage := make(map[string]int)
	//list := []chemical{equations["FUEL"].result}
	//for onlyOREBuilder(list) {
	//	fmt.Println("\n---------")
	//	fmt.Println("LIST:", list)
	//	goal := list[0]
	//	list = list[1:]
	//	fmt.Println("GOAL:", goal)
	//
	//	//if _, found := baseChemical[goal.name]; found {
	//	//	//solution := findChemicals(equations, goal)
	//	//	//fmt.Println("CS:", solution)
	//	//	baseChemical[goal.name] += goal.quantity
	//	//	fmt.Println("BASE found for ", goal.name, "with", equations[goal.name].chemicals[0], ", now", baseChemical[goal.name])
	//	//	continue
	//	//}
	//
	//	// Check how much we can get from storage.
	//	inStorage := storage[goal.name]
	//	if inStorage > 0 {
	//		tmp := goal.quantity
	//		if goal.quantity >= inStorage {
	//			goal.quantity -= inStorage
	//			storage[goal.name] -= tmp
	//			fmt.Println("Taking", inStorage, "from storage, left is", storage[goal.name])
	//		} else {
	//			goal.quantity = 0
	//			storage[goal.name] -= tmp
	//			fmt.Println("Taking", inStorage, "from storage, left is", storage[goal.name])
	//			continue
	//		}
	//	}
	//
	//	solution := findChemicals(equations, goal)
	//	fmt.Println(solution)
	//	if solution.result.quantity > goal.quantity {
	//		fmt.Println("Adding to storage for", goal.name, " amount=", solution.result.quantity-goal.quantity)
	//		storage[goal.name] += solution.result.quantity - goal.quantity
	//	}
	//
	//	if _, found := baseChemical[goal.name]; found {
	//		fmt.Println("Adding BASE", solution)
	//		list = append(list, solution.result)
	//	} else {
	//		fmt.Println("Adding", solution)
	//		list = append(list, solution.chemicals...)
	//	}
	//
	//	fmt.Print("<Enter>")
	//	bufio.NewReader(os.Stdin).ReadLine()
	//}
	//fmt.Println("\n\n\nSOLUTION")
	//ore := 0
	//for key, needed := range baseChemical {
	//	e := equations[key]
	//	q := e.result.quantity
	//	factor := int(math.Ceil(float64(needed) / float64(q)))
	//	o := e.chemicals[0].quantity * factor
	//	ore += o
	//	fmt.Println(key, factor, o)
	//}
	//fmt.Println("ORE:", ore, " delta=", 180697-ore)
}
