package main

import (
	"fmt"
	"math/rand/v2"
)

//go:wasmexport main
func main() {
	n := rand.IntN(99) + 1

	fmt.Println("Guess number from 1 to 100")

	var guess int
	for guess != n {
		_, err := fmt.Scanln(&guess)
		if err != nil {
			fmt.Println("Please enter a number")
		}

		if guess < n {
			fmt.Println("Too low")
		} else if guess > n {
			fmt.Println("Too high")
		}
	}

	fmt.Println("Correct!")
}
