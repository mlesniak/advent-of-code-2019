package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

	d := newDeck(10007)
	commands := load()
	fmt.Println(commands)
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
