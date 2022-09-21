package main

import "fmt"

func main() {

	fmt.Println(calcTotal(4000000, "even"))
}

// calcTotal calculates the total of a fibonacci sequence up to a max, with a specified rule for each value
func calcTotal(max int, rule string) int {
	fibs := fibLoop(max)
	total := 0

	switch rule {
	case "even":
		for _, fib := range fibs {
			if fib%2 == 0 {
				total += fib
			}
		}
	}

	return total

}

// fibLoop loops through all the fibonacci numbers in a sequence up to a maximum and appends them to an slice
func fibLoop(max int) []int {
	fibs := []int{1, 1}
	i1, i2, next := 1, 1, 0

	for {
		next = i1 + i2
		if next > max {
			break
		}

		i1 = i2
		i2 = next
		fibs = append(fibs, next)
	}

	return fibs
}
