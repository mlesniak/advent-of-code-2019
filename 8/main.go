package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

type layer [][]int

type image struct {
	width  int
	height int
	layers []layer
}

func main() {
	data := load()
	fmt.Println(data[:10])

}

func load() []int {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	chars := string(bytes)
	var memory []int
	for i := range chars {
		i, err := strconv.Atoi(string(chars[i]))
		if err != nil {
			panic(err)
		}
		memory = append(memory, i)
	}
	return memory
}
