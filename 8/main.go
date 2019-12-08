package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

type layer [][]int

type image struct {
	width  int
	height int
	layers []layer
}

func main() {
	width := 25
	height := 6
	image := image{width: width, height: height}
	data := load()

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

	// Would be easier if we had not split up the layers into two dimensions. ¯\_(ツ)_/¯
	minZeroes := math.MaxInt64
	ones := 0
	twos := 0
	for _, layer := range image.layers {
		z := 0
		o := 0
		t := 0
		for _, row := range layer {
			for _, value := range row {
				switch value {
				case 0:
					z++
				case 1:
					o++
				case 2:
					t++
				}
			}
		}
		if z < minZeroes {
			ones = o
			twos = t
			minZeroes = z
		}
	}

	fmt.Println(ones * twos)
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
