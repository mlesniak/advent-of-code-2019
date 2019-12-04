package main

import (
	"log"
	"strconv"
)

func main() {
	const start = 134792
	const end = 675810

	valid := 0
	for num := start; num <= end; num++ {
		if isValidPassword(strconv.Itoa(num)) {
			log.Printf("%d VALID\n", num)
			valid++
		} else {
			//log.Printf("%d NOTOK\n", num)
		}
	}

	log.Printf("Valid passwords: %d\n", valid)

	// Test cases from webpage.
	//fmt.Println(isValidPassword("111111"))
	//fmt.Println(isValidPassword("123789"))
	//fmt.Println(isValidPassword("223450"))
	//fmt.Println(isValidPassword("135578"))
	//fmt.Println(string(135578))
}

func isValidPassword(password string) bool {
	if len(password) != 6 {
		return false
	}

	// Check for increase or equal on subsequent digits.
	for idx := 0; idx < len(password)-1; idx++ {
		a, _ := strconv.Atoi(string(password[idx]))
		b, _ := strconv.Atoi(string(password[idx+1]))
		if a > b || a != 0 && b == 0 {
			return false
		}
	}

	// Check for at least one pair of double digits.
	foundDouble := false
	for idx := 0; idx < len(password)-1; idx++ {
		if password[idx] == password[idx+1] {
			foundDouble = true
			break
		}
	}
	if !foundDouble {
		return false
	}

	// All conditions fulfilled.
	return true
}
