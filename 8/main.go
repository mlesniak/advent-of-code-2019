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
	width := 3
	height := 2
	data := load()

	image := image{width: width, height: height}

	// Read image parts
	layerSize := width * height
	for layerIndex := 0; layerIndex < len(data)/layerSize; layerIndex++ {
		layerData := data[layerIndex*layerSize : layerIndex*layerSize+layerSize]
		layer := make([][]int, height)
		for row := 0; row < height; row++ {
			layer[row] = make([]int, width)
			for col := 0; col < width; col++ {
				idx := row*width + col
				layer[row][col] = layerData[idx]
			}
		}
		image.layers = append(image.layers, layer)
	}

	fmt.Println(image)
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
