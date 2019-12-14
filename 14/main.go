package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//10 ORE => 10 A
//1 ORE => 1 B
//7 A, 1 B => 1 C
//7 A, 1 C => 1 D
//7 A, 1 D => 1 E
//7 A, 1 E => 1 FUEL

type chemical struct {
	name     string
	quantity int
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
