package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func main() {
	maxBoost := math.MinInt64
	combinations := permutations([]int{5, 6, 7, 8, 9})
	for _, combination := range combinations {
		fmt.Println("Combination", combination)

		// Easier to reason about and play with using explicit loop unrolling in this case.
		a := load()
		aIn := make(chan int, 10)
		aOut := make(chan int, 10)
		go func() {
			compute("a", a, aIn, aOut)
		}()

		b := load()
		bIn := aOut
		bOut := make(chan int, 10)
		go func() {
			compute("b", b, bIn, bOut)
		}()

		c := load()
		cIn := bOut
		cOut := make(chan int, 10)
		go func() {
			compute("c", c, cIn, cOut)
		}()

		d := load()
		dIn := cOut
		dOut := make(chan int, 10)
		go func() {
			compute("d", d, dIn, dOut)
		}()

		e := load()
		eIn := dOut
		eOut := aIn

		// Initialize phase setting.
		aIn <- combination[0]
		bIn <- combination[1]
		cIn <- combination[2]
		dIn <- combination[3]
		eIn <- combination[4]

		// Start run
		aIn <- 0
		compute("e", e, eIn, eOut)

		n := <-eOut
		if n > maxBoost {
			maxBoost = n
		}
	}

	fmt.Println("MaxBoos", maxBoost)
}

// See https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func showResult(memory []int) {
	//fmt.Println(memory)
	fmt.Println(memory[0])
}

func compute(name string, memory []int, in chan int, out chan int) {
	for ip := 0; ip < len(memory); {
		mem := memory[ip]
		opcode := mem % 100
		r1 := mem / 100 % 10
		r2 := mem / 1000 % 10
		//r3 := memory[ip] / 10000 % 10

		switch opcode {
		case 1:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			var m3 int
			m3 = memory[ip+3]
			memory[m3] = m1 + m2
			ip += 4
		case 2:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			var m3 int
			m3 = memory[ip+3]
			memory[m3] = m1 * m2
			ip += 4
		case 3:
			var num int
			num = <-in
			memory[memory[ip+1]] = num
			ip += 2
		case 4:
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			//_, _ = fmt.Fprintf(out, "%d ", m1)
			out <- m1
			ip += 2
		case 5:
			// Jump if true.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if m1 != 0 {
				ip = m2
			} else {
				ip += 3
			}
		case 6:
			// Jump if false.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			if m1 == 0 {
				ip = m2
			} else {
				ip += 3
			}
		case 7:
			// Less than.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			var m3 int
			m3 = memory[ip+3]
			if m1 < m2 {
				memory[m3] = 1
			} else {
				memory[m3] = 0
			}
			ip += 4
		case 8:
			// Equals.
			var m1 int
			if r1 == 0 {
				m1 = memory[memory[ip+1]]
			}
			if r1 == 1 {
				m1 = memory[ip+1]
			}
			var m2 int
			if r2 == 0 {
				m2 = memory[memory[ip+2]]
			}
			if r2 == 1 {
				m2 = memory[ip+2]
			}
			var m3 int
			m3 = memory[ip+3]
			if m1 == m2 {
				memory[m3] = 1
			} else {
				memory[m3] = 0
			}
			ip += 4
		case 99:
			return
		default:
			panic("Unknown opcode:" + strconv.Itoa(opcode))
		}
	}
}

func load() []int {
	// Import.
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), ",")
	//fmt.Println(lines)
	var memory []int
	for _, val := range lines {
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		memory = append(memory, i)
	}
	return memory
}
