package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

// TODO Is a typedef sufficient?
type deck struct {
	cards []int
}

func (d deck) String() string {
	return fmt.Sprintf("CARDS: %v", d.cards)
}

func newDeck(size int) deck {
	nd := make([]int, size)
	for i := 0; i < size; i++ {
		nd[i] = i
	}

	return deck{cards: nd}
}

func (d deck) increment(n int) deck {
	nd := make([]int, len(d.cards))
	idx := 0
	for i := 0; i < len(d.cards); i++ {
		nd[idx] = d.cards[i]
		idx = (idx + n) % len(d.cards)
	}

	return deck{cards: nd}
}

func (d deck) cut(n int) deck {
	if n < 0 {
		n = len(d.cards) - n*-1
	}

	nd := make([]int, 0)
	nd = append(nd, d.cards[n:]...)
	nd = append(nd, d.cards[:n]...)

	return deck{cards: nd}
}

func (d deck) deal() deck {
	nd := make([]int, len(d.cards))
	for index, value := range d.cards {
		nd[len(nd)-1-index] = value
	}

	return deck{cards: nd}
}

func main() {
	//d := newDeck(10)
	//d = d.increment(7)
	//d = d.deal()
	//d = d.deal()
	//d = d.cut(6).increment(7).deal()
	//d = d.increment(7).increment(9).cut(-2)
	//fmt.Println(d)

	// First went with a totally wrong approach of using partial endings (stored in a separate branch), which took
	// me a couple of hours to realize that this won't work. Then took a look at the reddit solution thread to
	// find the theory behind the solution (much more mathy than programming).
	//
	// I'm not sure I understood all of the involved math correctly (linear functions in combination with
	// inverse modulos and relative prime numbers); my last number theory lecture was quite offset while ago.
	//
	// The logic has been aggreagted by reddit and other sources (not so much about programming, which is
	// this year's goal in AdventOfCode for me) -- hence I'm ok with it :-).
	//
	// See, e.g. https://www.reddit.com/r/adventofcode/comments/ee0rqi/2019_day_22_solutions/

	taskN := int64(119315717514047)
	count := int64(101741582076661)
	position := int64(2020)

	// Linear equation stuff...
	// Numbers are getting really big, took some time to figure this out...
	a := big.NewInt(0)
	b := big.NewInt(1)

	commands := load()
	for _, command := range commands {
		if command == "deal into new stack" {
			// increment *= -1
			b.Mul(b, big.NewInt(-1))
			//offset += increment
			a.Add(a, b)
			continue
		}

		cut := int64(0)
		n, _ := fmt.Sscanf(command, "cut %d\n", &cut)
		if n > 0 {
			//offset = offset + (increment * cut)
			a.Add(a, big.NewInt(0).Mul(b, big.NewInt(cut)))
			continue
		}

		inc := 0
		n, _ = fmt.Sscanf(command, "deal with increment %d\n", &inc)
		if n > 0 {
			// We need exponentiation with modulo, hence bigInts for this case.
			nInt := big.NewInt(taskN)
			bi := big.NewInt(0).Exp(big.NewInt(int64(inc)), big.NewInt(0).Sub(nInt, big.NewInt(2)), nInt)
			//increment = increment * bi.Int64()
			b.Mul(b, bi)
			continue
		}
	}

	bigN := big.NewInt(taskN)
	na := big.NewInt(0).Exp(b, big.NewInt(count), bigN)

	nb := big.NewInt(0).Exp(b, big.NewInt(count), bigN)
	nb.Sub(big.NewInt(1), nb)
	// Number theory stuff.
	invmod := big.NewInt(0).Exp(big.NewInt(0).Sub(big.NewInt(1), b), big.NewInt(0).Sub(bigN, big.NewInt(2)), bigN)
	nb.Mul(nb.Mul(nb, a), invmod)

	result := big.NewInt(0).Mul(big.NewInt(position), na)
	result.Add(result, nb)
	result.Mod(result, bigN)

	fmt.Println(result)
}

func part1() {
	d := newDeck(10007)
	commands := load()
	for _, command := range commands {
		if command == "deal into new stack" {
			fmt.Println("Dealing new stack")
			d = d.deal()
			continue
		}

		cut := 0
		n, _ := fmt.Sscanf(command, "cut %d\n", &cut)
		if n > 0 {
			fmt.Println("Cutting", cut)
			d = d.cut(cut)
			continue
		}

		inc := 0
		n, _ = fmt.Sscanf(command, "deal with increment %d\n", &inc)
		if n > 0 {
			fmt.Println("Incrementing", inc)
			d = d.increment(inc)
			continue
		}
	}
	for key, value := range d.cards {
		if value == 2019 {
			fmt.Println(key)
			break
		}
	}
}

// TODO Parse into commands
func load() []string {
	bytes, _ := ioutil.ReadFile("input.txt")
	lines := strings.Split(string(bytes), "\n")

	return lines
}

func wait() {
	fmt.Print("<ENTER>")
	bufio.NewReader(os.Stdin).ReadLine()
}
